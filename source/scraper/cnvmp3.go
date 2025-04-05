package scraper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"fyne.io/fyne/v2"
)

type cnvmp3Downloader struct {
	referer          string
	token            string
	downloadVideoURL string
}

func newCnvmp3Downloader() *cnvmp3Downloader {
	rsp, err := http.Get(`https://cnvmp3.com/`)
	if err != nil {
		fyne.LogError("failed to obtain cvnmp3 download video url", err)
	}
	defer rsp.Body.Close()

	content, err := io.ReadAll(rsp.Body)
	if err != nil {
		fyne.LogError("failed to decode cvnmp3 source", err)
	}

	// Scrape referer
	referer := regexp.
		MustCompile(`<link rel="canonical" href="(.+)">`).
		FindStringSubmatch(string(content))[1]

	// Scrape download token.
	downloadVideoToken := regexp.
		MustCompile(`data.token = \"(.+)\";`).
		FindStringSubmatch(string(content))[1]

	// Scrape download URL.
	downloadVidelURL := regexp.
		MustCompile(`function downloadVideo\(.+\) \{.+\n.+fetch\('(.+)', \{`).
		FindStringSubmatch(string(content))[1]

	return &cnvmp3Downloader{referer, downloadVideoToken, `https://cnvmp3.com/` + downloadVidelURL}
}

func (d *cnvmp3Downloader) Download(video *Result) (io.ReadCloser, error) {
	// Skip the database part, ignore the serverside caching.
	if err := d.getVideoData(video); err != nil {
		return nil, err
	}
	filelink, err := d.getVideoDownloadLink(video)
	if err != nil {
		return nil, err
	}

	// Download the music file
	req, err := http.NewRequest(http.MethodGet, filelink, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("host", "apiv17dlp.cnvmp3.me")
	req.Header.Set("referer", d.referer)
	musicResp, err := http.DefaultClient.Do(req)
	return musicResp.Body, err
}

func (d *cnvmp3Downloader) getVideoData(video *Result) error {
	type GetVideoDataRequest struct {
		Token string `json:"token"`
		URL   string `json:"url"`
	}

	type GetVideoDataResponse struct {
		Success bool   `json:"success"`
		Title   string `json:"title"`
	}

	// Prepare the request.
	const endpoint = `https://cnvmp3.com/get_video_data.php`
	request := GetVideoDataRequest{Token: d.token, URL: `https://www.youtube.com/watch?` + url.Values{"v": {video.ID}}.Encode()}
	requestData, err := json.Marshal(request)
	if err != nil {
		return err
	}

	// Send the request.
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(requestData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Parse the response.
	response := GetVideoDataResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}
	if !response.Success {
		return fmt.Errorf("failed to get video data")
	}
	return nil
}

func (d *cnvmp3Downloader) getVideoDownloadLink(video *Result) (string, error) {
	type DownloadVideoRequest struct {
		URL         string `json:"url"`
		Quality     int64  `json:"quality"`
		Title       string `json:"title"`
		FormatValue int64  `json:"formatValue"`
	}

	type DownloadVideoResponse struct {
		Success      bool   `json:"success"`
		DownloadLink string `json:"download_link"`
	}

	// Prepare the request.
	request := DownloadVideoRequest{
		URL:         `https://www.youtube.com/watch?` + url.Values{"v": {video.ID}}.Encode(),
		Quality:     0,
		Title:       video.Title,
		FormatValue: 1,
	}
	requestData, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	// Send the request.
	req, err := http.NewRequest(http.MethodPost, d.downloadVideoURL, bytes.NewBuffer(requestData))
	req.Header.Set("referer", d.referer)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse the response.
	response := DownloadVideoResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	if !response.Success {
		return "", fmt.Errorf("failed to get download video link")
	}

	paramCutOff := strings.Index(response.DownloadLink, "=") + 1
	response.DownloadLink = response.DownloadLink[:paramCutOff] + url.QueryEscape(response.DownloadLink[paramCutOff:])
	return response.DownloadLink, nil
}

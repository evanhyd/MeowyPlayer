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

var cnvmp3DownloadVideoURL string

func init() {
	rsp, err := http.Get(`https://cnvmp3.com/`)
	if err != nil {
		fyne.LogError("failed to obtain cvnmp3 download video url", err)
	}
	defer rsp.Body.Close()

	content, err := io.ReadAll(rsp.Body)
	if err != nil {
		fyne.LogError("failed to decode cvnmp3 source", err)
	}

	exp := regexp.MustCompile(`function downloadVideo\(.+\) \{.+\n.+fetch\('(.+)', \{`)
	cnvmp3DownloadVideoURL = exp.FindStringSubmatch(string(content))[1]
}

type cnvmp3Downloader struct {
	downloadVideoURL string
}

func newCnvmp3Downloader() *cnvmp3Downloader {
	return &cnvmp3Downloader{`https://cnvmp3.com/` + cnvmp3DownloadVideoURL}
}

func (d *cnvmp3Downloader) Download(video *Result) (io.ReadCloser, error) {
	// Skip the database part, ignore the serverside caching.
	if err := d.getVideoData(video); err != nil {
		return nil, err
	}
	filelink, err := d.downloadVideo(video)
	if err != nil {
		return nil, err
	}
	paramCutOff := strings.Index(filelink, "?") + 1
	filelink = filelink[:paramCutOff] + url.PathEscape(filelink[paramCutOff:])

	// Download the music file
	downloadRequest, err := http.NewRequest(http.MethodGet, filelink, nil)
	if err != nil {
		return nil, err
	}
	downloadRequest.Header.Set("host", "apiv17dlp.cnvmp3.me")
	downloadRequest.Header.Set("referer", "https://cnvmp3.com/")
	musicResp, err := http.DefaultClient.Do(downloadRequest)
	return musicResp.Body, err
}

func (d *cnvmp3Downloader) getVideoData(video *Result) error {
	type GetVideoDataRequest struct {
		URL string `json:"url"`
	}

	type GetVideoDataResposne struct {
		Success bool   `json:"success"`
		Title   string `json:"title"`
	}

	// Prepare the request.
	const endpoint = `https://cnvmp3.com/get_video_data.php`
	request := GetVideoDataRequest{URL: `https://www.youtube.com/watch?` + url.Values{"v": {video.ID}}.Encode()}
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
	response := GetVideoDataResposne{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}
	if !response.Success {
		return fmt.Errorf("failed to get video data")
	}
	return nil
}

func (d *cnvmp3Downloader) downloadVideo(video *Result) (string, error) {
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
	req.Header.Set("referer", `https://cnvmp3.com`)
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
	return response.DownloadLink, nil
}

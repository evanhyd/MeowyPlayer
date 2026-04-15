package scrapers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type cnvmp3Downloader struct {
	referer          string
	token            string
	downloadVideoURL string
}

func NewCnvmp3Downloader() *cnvmp3Downloader {
	rsp, err := http.Get(`https://cnvmp3.com/`)
	if err != nil {
		slog.Error("failed to obtain cvnmp3 download video url", "error", err)
		return &cnvmp3Downloader{}
	}
	defer rsp.Body.Close()

	content, err := io.ReadAll(rsp.Body)
	if err != nil {
		slog.Error("failed to decode cvnmp3 source", "error", err)
		return &cnvmp3Downloader{}
	}

	// Scrape referer.
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

func (d *cnvmp3Downloader) Download(ctx context.Context, video Result) (io.ReadCloser, error) {
	if err := d.getVideoData(ctx, &video); err != nil {
		return nil, err
	}
	filelink, err := d.getVideoDownloadLink(ctx, &video)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, filelink, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("host", "apiv17dlp.cnvmp3.me")
	req.Header.Set("referer", d.referer)

	musicResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err // Returns context.Canceled if the user aborted
	}

	return musicResp.Body, nil
}

func (d *cnvmp3Downloader) getVideoData(ctx context.Context, video *Result) error {
	type GetVideoDataRequest struct {
		Token string `json:"token"`
		URL   string `json:"url"`
	}

	type GetVideoDataResponse struct {
		Success bool   `json:"success"`
		Title   string `json:"title"`
	}

	const endpoint = `https://cnvmp3.com/get_video_data.php`
	request := GetVideoDataRequest{Token: d.token, URL: `https://www.youtube.com/watch?` + url.Values{"v": {video.ID}}.Encode()}
	requestData, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(requestData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	response := GetVideoDataResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}
	if !response.Success {
		return fmt.Errorf("failed to get video data")
	}
	return nil
}

func (d *cnvmp3Downloader) getVideoDownloadLink(ctx context.Context, video *Result) (string, error) {
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

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, d.downloadVideoURL, bytes.NewBuffer(requestData))
	if err != nil {
		return "", err
	}
	req.Header.Set("referer", d.referer)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

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

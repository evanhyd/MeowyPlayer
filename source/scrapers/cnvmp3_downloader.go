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

var _ MusicDownloader = (*cnvmp3Downloader)(nil)

type cnvmp3Downloader struct {
	referer                 string
	token                   string
	downloadVideoScriptPath string
}

func NewCnvmp3Downloader() *cnvmp3Downloader {
	rsp, err := http.Get(`https://cnvmp3.com/`)
	if err != nil {
		slog.Error("failed to obtain cvnmp3 download video url", "error", err)
		return nil
	}
	defer rsp.Body.Close()

	content, _ := io.ReadAll(rsp.Body)
	str := string(content)

	referer := regexp.MustCompile(`<link rel="canonical" href="(.+)">`).FindStringSubmatch(str)[1]
	token := regexp.MustCompile(`data\.token = "(.+)";`).FindStringSubmatch(str)[1]
	path := regexp.MustCompile(`function downloadVideo\(.+\) \{.+\n.+fetch\('(.+)', \{`).FindStringSubmatch(str)[1]

	return &cnvmp3Downloader{
		referer:                 referer,
		token:                   token,
		downloadVideoScriptPath: `https://cnvmp3.com/` + path,
	}
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
		return nil, err
	}

	// Prevents saving JSON error responses as MP3 files
	if musicResp.StatusCode != http.StatusOK || strings.Contains(musicResp.Header.Get("Content-Type"), "json") {
		defer musicResp.Body.Close()
		body, _ := io.ReadAll(musicResp.Body)
		return nil, fmt.Errorf("download failed (%d): %s", musicResp.StatusCode, string(body))
	}

	return musicResp.Body, nil
}

func (d *cnvmp3Downloader) getVideoData(ctx context.Context, video *Result) error {
	data, _ := json.Marshal(map[string]string{
		"token": d.token,
		"url":   "https://www.youtube.com/watch?v=" + video.ID,
	})

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "https://cnvmp3.com/get_video_data.php", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var res struct {
		Success bool   `json:"success"`
		Title   string `json:"title"`
		Error   string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}
	if res.Error != "" {
		return fmt.Errorf("video data error: %s", res.Error)
	}

	if video.Title == "" {
		video.Title = res.Title
	}
	return nil
}

func (d *cnvmp3Downloader) getVideoDownloadLink(ctx context.Context, video *Result) (string, error) {
	data, _ := json.Marshal(map[string]any{
		"url":         "https://www.youtube.com/watch?v=" + video.ID,
		"quality":     0,
		"title":       video.Title,
		"formatValue": 1,
	})

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, d.downloadVideoScriptPath, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("referer", d.referer)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res struct {
		Error        string `json:"error"`
		DownloadLink string `json:"download_link"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}
	if res.Error != "" || res.DownloadLink == "" {
		return "", fmt.Errorf("download link error: %s", res.Error)
	}

	paramCutOff := strings.Index(res.DownloadLink, "=") + 1
	res.DownloadLink = res.DownloadLink[:paramCutOff] + url.PathEscape(res.DownloadLink[paramCutOff:])
	return res.DownloadLink, nil
}

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
	"time"
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
		slog.Error("failed to obtain cnvmp3 download video url", "error", err)
		return nil
	}
	defer rsp.Body.Close()

	content, _ := io.ReadAll(rsp.Body)
	str := string(content)

	referer := regexp.MustCompile(`<link rel="canonical" href="([^"]+)">`).FindStringSubmatch(str)[1]
	token := regexp.MustCompile(`data\.token\s*=\s*"([^"]+)";`).FindStringSubmatch(str)[1]
	path := regexp.MustCompile(`(?s)function downloadVideo.*?fetch\('([^']+)'`).FindStringSubmatch(str)[1]

	return &cnvmp3Downloader{
		referer:                 referer,
		token:                   token,
		downloadVideoScriptPath: `https://cnvmp3.com/` + strings.TrimPrefix(path, "/"),
	}
}

func (d *cnvmp3Downloader) Download(ctx context.Context, video Result) (io.ReadCloser, error) {
	filelink, err := d.getVideoDownloadLink(ctx, &video)
	if err != nil {
		if err := d.getVideoData(ctx, &video); err != nil {
			return nil, err
		}
		filelink, err = d.requestConversion(ctx, &video)
		if err != nil {
			return nil, err
		}
	}

	for attempt := 0; attempt < 10; attempt++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, filelink, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("referer", d.referer)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

		musicResp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}

		contentType := musicResp.Header.Get("Content-Type")
		if musicResp.StatusCode == http.StatusOK && !strings.Contains(contentType, "json") {
			return musicResp.Body, nil
		}

		body, _ := io.ReadAll(musicResp.Body)
		musicResp.Body.Close()

		if strings.Contains(string(body), "File not found") && attempt < 9 {
			time.Sleep(1500 * time.Millisecond)
			continue
		}

		return nil, fmt.Errorf("download failed (%d): %s", musicResp.StatusCode, string(body))
	}

	return nil, fmt.Errorf("download timed out waiting for server")
}

func (d *cnvmp3Downloader) getVideoData(ctx context.Context, video *Result) error {
	data, _ := json.Marshal(map[string]string{
		"token": d.token,
		"url":   "https://www.youtube.com/watch?v=" + video.ID,
	})

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "https://cnvmp3.com/get_video_data.php", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("referer", d.referer)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var res struct {
		Title string `json:"title"`
		Error string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("download failed (%d): %s", resp.StatusCode, string(body))
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
		"youtube_id":  video.ID,
		"quality":     0, // 0 = 320kbps
		"formatValue": 1, // 1 = MP3
	})

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "https://cnvmp3.com/check_database.php", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("referer", d.referer)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res struct {
		Success bool `json:"success"`
		Data    struct {
			ServerPath string `json:"server_path"`
			Title      string `json:"title"`
		} `json:"data"`
		Error any `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	if !res.Success || res.Data.ServerPath == "" {
		return "", fmt.Errorf("not in database")
	}

	if video.Title == "" && res.Data.Title != "" {
		video.Title = res.Data.Title
	}

	return formatDownloadURL(res.Data.ServerPath), nil
}

func (d *cnvmp3Downloader) requestConversion(ctx context.Context, video *Result) (string, error) {
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
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("conversion script failed (%d): %s", resp.StatusCode, string(body))
	}
	if res.Error != "" || res.DownloadLink == "" {
		return "", fmt.Errorf("conversion error: %s", res.Error)
	}

	return formatDownloadURL(res.DownloadLink), nil
}

func formatDownloadURL(rawURL string) string {
	idx := strings.Index(rawURL, "?file=")
	if idx == -1 {
		return rawURL
	}

	base := rawURL[:idx+len("?file=")]
	filePath := rawURL[idx+len("?file="):]

	parts := strings.Split(filePath, "/")
	for i, part := range parts {
		escaped := url.QueryEscape(part)
		escaped = strings.ReplaceAll(escaped, "+", "%20")
		parts[i] = escaped
	}

	return base + strings.Join(parts, "/")
}

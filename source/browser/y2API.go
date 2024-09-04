package browser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	kAPIProvider = `https://download.y2api.com`
)

// https://download.y2api.com/
type y2APIDownloader struct {
}

func newY2APIDownloader() *y2APIDownloader {
	return &y2APIDownloader{}
}

// extract the converter link from the widget
func (d *y2APIDownloader) getConverter(video *Result) (string, string, error) {

	type y2APIMp4 struct {
		Quality string `json:"quality"`
		Itag    int64  `json:"itag"`
		EXT     string `json:"ext"`
		Codec   string `json:"codec"`
		Token   string `json:"token"`
	}

	type y2APIVideo struct {
		Mp4 []y2APIMp4 `json:"mp4"`
	}

	type y2APIMp3 struct {
		Quality int64  `json:"quality"`
		Codec   string `json:"codec"`
		EXT     string `json:"ext"`
		Token   string `json:"token"`
	}

	type y2APIAudio struct {
		Mp3 []y2APIMp3 `json:"mp3"`
	}

	type y2APIFormats struct {
		Audio y2APIAudio `json:"audio"`
		Video y2APIVideo `json:"video"`
	}

	type y2APIWidgetResponse struct {
		VideoID       string       `json:"videoId"`
		Title         string       `json:"title"`
		Duration      string       `json:"duration"`
		HumanDuration string       `json:"humanDuration"`
		Formats       y2APIFormats `json:"formats"`
	}

	//get converter
	reqs, err := http.NewRequest(http.MethodGet, `https://rr-01-bucket.cdn1313.net/api/v4/info/`+video.ID, nil)
	if err != nil {
		return "", "", err
	}
	reqs.Header.Set("origin", kAPIProvider)
	reqs.Header.Set("referer", kAPIProvider)
	resp, err := http.DefaultClient.Do(reqs)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	//parse converter
	widgetResp := y2APIWidgetResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&widgetResp); err != nil {
		return "", "", err
	}
	if len(widgetResp.Formats.Audio.Mp3) == 0 {
		return "", "", fmt.Errorf("video is too long or unavailable")
	}

	return widgetResp.Formats.Audio.Mp3[0].Token, resp.Header.Get("authorization"), nil
}

// request the conversion and return the job ID
func (d *y2APIDownloader) getJob(token string, authorization string) (string, error) {
	type y2APIConverterPayload struct {
		Token string `json:"token"`
	}
	type y2APIConverterResponse struct {
		ID      string `json:"id"`
		Message string `json:"message"`
	}

	//request conversion
	payload := y2APIConverterPayload{Token: token}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	reqs, err := http.NewRequest(http.MethodPost, `https://rr-01-bucket.cdn1313.net/api/v4/convert`, bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	reqs.Header.Set("authorization", authorization)
	reqs.Header.Set("origin", kAPIProvider)
	reqs.Header.Set("referer", kAPIProvider)
	resp, err := http.DefaultClient.Do(reqs)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	//parse converter response
	converterResponse := y2APIConverterResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&converterResponse); err != nil {
		return "", err
	}
	return converterResponse.ID, nil
}

// monitor the conversion job status and return the download link
func (d *y2APIDownloader) getDownload(authorization string, jobID string) (string, error) {
	type y2APIJobResponse struct {
		ID       string  `json:"id"`
		Status   string  `json:"status"`
		Progress float32 `json:"progress"`
		VideoID  string  `json:"videoId"`
		Title    string  `json:"title"`
		EXT      string  `json:"ext"`
		Quality  int64   `json:"quality"`
		Download string  `json:"download"`
	}

	reqs, err := http.NewRequest(http.MethodGet, `https://rr-01-bucket.cdn1313.net/api/v4/status/`+jobID, nil)
	if err != nil {
		return "", err
	}
	reqs.Header.Set("authorization", authorization)
	reqs.Header.Set("origin", kAPIProvider)
	reqs.Header.Set("referer", kAPIProvider)

	//parse download response
	for {
		resp, err := http.DefaultClient.Do(reqs)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		jobResp := y2APIJobResponse{}
		if err := json.NewDecoder(resp.Body).Decode(&jobResp); err != nil {
			return "", err
		}

		if jobResp.Status == "completed" {
			return jobResp.Download, nil
		}
		if jobResp.Status == "active" {
			time.Sleep(1 * time.Second)
			continue
		}
		return "", fmt.Errorf("unknown error: %v", jobResp)
	}
}

func (d *y2APIDownloader) Download(video *Result) (io.ReadCloser, error) {
	token, authorization, err := d.getConverter(video)
	if err != nil {
		return nil, err
	}

	jobID, err := d.getJob(token, authorization)
	if err != nil {
		return nil, err
	}

	downloadURL, err := d.getDownload(authorization, jobID)
	if err != nil {
		return nil, err
	}

	//request download
	resp, err := http.Get(downloadURL)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

package browser

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
)

// https://app.quicktype.io/
type y2mateConverterResponse struct {
	Status   string `json:"status"`
	Mess     string `json:"mess"`
	CStatus  string `json:"c_status"`
	Vid      string `json:"vid"`
	Title    string `json:"title"`
	Ftype    string `json:"ftype"`
	Fquality string `json:"fquality"`
	Dlink    string `json:"dlink"`
}

type y2MateDownloader struct {
	keyRegex *regexp.Regexp
}

func newY2MateDownloader() *y2MateDownloader {
	const keyPattern = `"f":"mp3","q":"128kbps","q_text":"MP3 - 128kbps","k":"([\w\/\\\+=]+)"`
	return &y2MateDownloader{regexp.MustCompile(keyPattern)}
}

func (d *y2MateDownloader) Download(video *Result) (io.ReadCloser, error) {
	key, err := d.parseConverterKey(video)
	if err != nil {
		return nil, err
	}
	return d.downloadContent(video, key)
}

func (d *y2MateDownloader) parseConverterKey(video *Result) (string, error) {
	const (
		converterUrl = `https://www.y2mate.com/mates/analyzeV2/ajax`
		youtubeUrl   = `https://www.youtube.com/watch?`
	)

	//request the content that contains converter key
	videoUrl := youtubeUrl + url.Values{"v": {video.VideoID}}.Encode()
	queryData := url.Values{"k_query": {videoUrl}, "k_page": {"home"}, "hl": {"en"}, "q_auto": {"1"}}
	resp, err := http.PostForm(converterUrl, queryData)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	//regex parse the key from the converter
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	matches := d.keyRegex.FindStringSubmatch(string(data))
	if matches == nil {
		return "", fmt.Errorf("couldn't get converter key: %v, %v", video.VideoID, video.Title)
	}
	return matches[1], nil
}

func (d *y2MateDownloader) downloadContent(video *Result, converterKey string) (io.ReadCloser, error) {
	const dbURL = `https://www.y2mate.com/mates/convertV2/index`

	//request for video -> mp3 conversion
	queryData := url.Values{"vid": {video.VideoID}, "k": {converterKey}}
	resp, err := http.PostForm(dbURL, queryData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//parse json response
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	converterResp := y2mateConverterResponse{}
	if err := json.Unmarshal(data, &converterResp); err != nil {
		return nil, err
	}

	//fetch music file
	if converterResp.CStatus == "FAILED" {
		return nil, fmt.Errorf("[Y2mate] failed to download %v, can not find the resource", video.VideoID)
	}
	resp, err = http.Get(converterResp.Dlink)
	return resp.Body, err
}

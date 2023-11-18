package downloader

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"meowyplayer.com/utility/network/fileformat"
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

func NewY2MateDownloader() *y2MateDownloader {
	const keyPattern = `"f":"mp3","q":"128kbps","q_text":"MP3 - 128kbps","k":"([\w\/\\\+=]+)"`
	return &y2MateDownloader{regexp.MustCompile(keyPattern)}
}

func (d *y2MateDownloader) Download(video *fileformat.VideoResult) ([]byte, error) {
	key, err := d.parseConverterKey(video)
	if err != nil {
		return nil, err
	}
	return d.downloadContent(video, key)
}

func (d *y2MateDownloader) parseConverterKey(video *fileformat.VideoResult) (string, error) {
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
		return "", fmt.Errorf("couldn't get converter key: %v", video.Title)
	}
	return matches[1], nil
}

func (d *y2MateDownloader) downloadContent(video *fileformat.VideoResult, converterKey string) ([]byte, error) {
	const dbURL = `https://www.y2mate.com/mates/convertV2/index`

	log.Printf("[Y2mate] fetching %v with key %v\n", video.Title, converterKey)

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
	resp, err = http.Get(converterResp.Dlink)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

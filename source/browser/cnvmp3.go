package browser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type cnvmp3Downloader struct {
}

func newCnvmp3Downloader() *cnvmp3Downloader {
	return &cnvmp3Downloader{}
}

func (d *cnvmp3Downloader) Download(video *Result) (io.ReadCloser, error) {
	type cnvmp3Payload struct {
		URL           string `json:"url"`
		DownloadMode  string `json:"downloadMode"`
		FilenameStyle string `json:"filenameStyle"`
		AudioBitrate  string `json:"audioBitrate"`
	}

	type cnvmp3Response struct {
		Status   string `json:"status"`
		URL      string `json:"url"`
		Filename string `json:"filename"`
	}

	payloadData, err := json.Marshal(cnvmp3Payload{URL: "https://www.youtube.com/watch?v=" + video.ID, DownloadMode: "audio", FilenameStyle: "pretty", AudioBitrate: "320"})
	if err != nil {
		return nil, err
	}

	convertReq, err := http.NewRequest(http.MethodPost, "https://cnvmp3.com/fetch.php", bytes.NewBuffer(payloadData))
	if err != nil {
		return nil, err
	}

	convertRsp, err := http.DefaultClient.Do(convertReq)
	if err != nil {
		return nil, err
	}
	defer convertRsp.Body.Close()

	response := cnvmp3Response{}
	if err := json.NewDecoder(convertRsp.Body).Decode(&response); err != nil {
		errData, _ := io.ReadAll(convertRsp.Body)
		defer convertRsp.Body.Close()
		return nil, fmt.Errorf("%s", errData)
	}

	if response.Status != "tunnel" { //"error"
		return nil, fmt.Errorf("%v", response.URL)
	}

	fileRsp, err := http.Get(response.URL)
	if err != nil {
		return nil, err
	}
	return fileRsp.Body, nil
}

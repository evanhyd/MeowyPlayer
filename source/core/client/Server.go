package client

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	"meowyplayer.com/core/resource"
	"meowyplayer.com/utility/uzip"
)

func sendRequest(
	method string, server string, queryType string, urlValues url.Values,
	username string, password string,
	contentType string, content io.Reader) (*http.Response, error) {

	//base url
	url, err := url.JoinPath(server, queryType)
	if err != nil {
		return nil, err
	}

	//url values
	if len(urlValues) > 0 {
		url += "?" + urlValues.Encode()
	}

	//create request
	req, err := http.NewRequest(method, url, content)
	if err != nil {
		return nil, err
	}

	//set auth
	req.SetBasicAuth(username, password)

	//set content type
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	return http.DefaultClient.Do(req)
}

func RequestList(server, username, password string) ([]resource.CollectionInfo, error) {
	resp, err := sendRequest("GET", server, "list", nil, username, password, "", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	infos := []resource.CollectionInfo{}
	err = json.NewDecoder(resp.Body).Decode(&infos)
	return infos, err
}

func RequestUpload(server, username, password string) error {
	//zip collection config
	zipData := bytes.Buffer{}
	if err := uzip.Compress(&zipData, resource.CollectionPath()); err != nil {
		return err
	}

	//prepare POST fields
	fieldBody := bytes.Buffer{}
	fieldWriter := multipart.NewWriter(&fieldBody)
	writeFields := func() error {
		defer fieldWriter.Close()
		fieldPart, err := fieldWriter.CreateFormFile("collection", Config().Name())
		if err != nil {
			return err
		}
		_, err = fieldPart.Write(zipData.Bytes())
		return err
	}
	if err := writeFields(); err != nil {
		return err
	}

	_, err := sendRequest("POST", server, "upload", nil, username, password, fieldWriter.FormDataContentType(), &fieldBody)
	return err
}

func RequestDownload(server, username, password string, info *resource.CollectionInfo) (<-chan float64, error) {
	resp, err := sendRequest("GET", server, "download", url.Values{"collection": {info.Title}}, username, password, "", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//read in zip format
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}

	//save to local
	if err := os.RemoveAll(resource.CollectionPath()); err != nil {
		return nil, err
	}
	if err := uzip.Extract(resource.CollectionPath(), reader); err != nil {
		return nil, err
	}
	if err := Manager().load(); err != nil {
		return nil, err
	}

	//sync music list
	return SyncCollection(), nil
}

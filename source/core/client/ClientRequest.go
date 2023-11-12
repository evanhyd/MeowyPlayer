package client

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"meowyplayer.com/core/resource"
	"meowyplayer.com/utility/uzip"
)

func (c *clientManager) ClientRequestList(account *resource.Account) ([]resource.CollectionInfo, error) {
	serverUrl, err := url.JoinPath(Config().ServerUrl, "list")
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(serverUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	infos := []resource.CollectionInfo{}
	err = json.NewDecoder(resp.Body).Decode(&infos)
	return infos, err
}

func (c *clientManager) ClientRequestUpload(account *resource.Account) error {
	serverUrl, err := url.JoinPath(Config().ServerUrl, "upload")
	if err != nil {
		return err
	}

	//zip files
	zipData := bytes.Buffer{}
	if err := uzip.Compress(&zipData, resource.CollectionPath()); err != nil {
		return err
	}

	//prepare POST fields
	fieldBody := bytes.Buffer{}
	fieldWriter := multipart.NewWriter(&fieldBody)

	writeFields := func() error {
		defer fieldWriter.Close()
		fieldPart, err := fieldWriter.CreateFormFile("collection", account.Name+".zip")
		if err != nil {
			return err
		}
		_, err = io.Copy(fieldPart, &zipData)
		return err
	}
	if err := writeFields(); err != nil {
		return err
	}

	//send post
	_, err = http.Post(serverUrl, fieldWriter.FormDataContentType(), &fieldBody)
	return err
}

func (c *clientManager) ClientRequestDownload(account *resource.Account, collectionInfo *resource.CollectionInfo) error {
	serverUrl, err := url.JoinPath(Config().ServerUrl, "download")
	if err != nil {
		return err
	}
	serverUrl += "?" + url.Values{"collection": {collectionInfo.Title}}.Encode()

	//download collection
	resp, err := http.Get(serverUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//unzip collection
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}

	c.accessLock.Lock()
	defer c.accessLock.Unlock()

	if err := uzip.Extract(resource.CollectionPath(), reader); err != nil {
		return err
	}
	return c.Load()
}

func (c *clientManager) ClientRequestRemove(account *resource.Account, collectionInfo *resource.CollectionInfo) error {
	// serverUrl, err := url.JoinPath(Config().ServerUrl, "remove")
	// if err != nil {
	// 	return err
	// }
	return nil
}

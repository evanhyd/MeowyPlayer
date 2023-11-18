package client

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	"meowyplayer.com/core/resource"
	"meowyplayer.com/utility/uzip"
)

func RequestList(account *resource.Account) ([]resource.CollectionInfo, error) {
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

func RequestUpload(account *resource.Account) error {
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

func downloadCollection(serverUrl string) error {
	resp, err := http.Get(serverUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}

	Manager().accessLock.Lock()
	defer Manager().accessLock.Unlock()

	if err := os.RemoveAll(resource.CollectionPath()); err != nil {
		return err
	}
	if err := uzip.Extract(resource.CollectionPath(), reader); err != nil {
		return err
	}
	return Manager().load()
}

func RequestDownload(account *resource.Account, collectionInfo *resource.CollectionInfo) error {
	serverUrl, err := url.JoinPath(Config().ServerUrl, "download")
	if err != nil {
		return err
	}
	serverUrl += "?" + url.Values{"collection": {collectionInfo.Title}}.Encode()
	if err := downloadCollection(serverUrl); err != nil {
		return err
	}
	if unsynced := SyncCollection(); unsynced != 0 {
		return fmt.Errorf("unable to sync %v music", unsynced)
	}
	return nil
}

func RequestRemove(account *resource.Account, collectionInfo *resource.CollectionInfo) error {
	// serverUrl, err := url.JoinPath(Config().ServerUrl, "remove")
	// if err != nil {
	// 	return err
	// }
	return nil
}

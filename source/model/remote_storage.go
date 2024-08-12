package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type remoteStorageConfig struct {
	ServerEndpoint string `json:"serverEndpoint"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	// Token          []byte `json:"token"`
}

type remoteStorage struct {
	configDir string
	config    remoteStorageConfig
}

var _ Storage = &remoteStorage{}

func newRemoteStorage() *remoteStorage {
	const kStorage = "storage"
	return &remoteStorage{
		configDir: filepath.Join(kStorage, "config.json"),
		config:    remoteStorageConfig{ServerEndpoint: "http://132.145.98.4/"},
	}
}

func (s *remoteStorage) initialize() error {
	data, err := os.ReadFile(s.configDir)

	//create default config
	if errors.Is(err, os.ErrNotExist) {
		return s.save()
	}

	//other error
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s.config)
}

func (s *remoteStorage) save() error {
	data, err := json.Marshal(s.config)
	if err != nil {
		return err
	}
	return os.WriteFile(s.configDir, data, 0600)
}

func (s *remoteStorage) setConfig(username string, password string) error {
	s.config.Username = username
	s.config.Password = password
	return s.save()
}

func (s *remoteStorage) sendPost(apiFunc string, username string, password string) (*http.Response, error) {
	url, err := url.JoinPath(s.config.ServerEndpoint, apiFunc)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(username, password)

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if rsp.StatusCode != http.StatusOK {
		errMsg, err := io.ReadAll(rsp.Body)
		if err != nil {
			return rsp, err
		}
		defer rsp.Body.Close()
		return nil, fmt.Errorf("%s", errMsg)
	}
	return rsp, nil
}

func (s *remoteStorage) authenticate(username string, password string) error {
	_, err := s.sendPost("login", username, password)
	return err
}

func (s *remoteStorage) register(username string, password string) error {
	_, err := s.sendPost("register", username, password)
	return err
}

// Get all the personal albums from the remote.
func (s *remoteStorage) getAllAlbums() ([]Album, error) {
	panic("not implemented yet")
}

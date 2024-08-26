package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"meowyplayer/util"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
)

type UserInfo struct {
	Username string
}

type networkClientConfig struct {
	ServerURL string `json:"serverURL"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	// Token          []byte `json:"token"`
}

type networkClient struct {
	configDir      string
	config         networkClientConfig
	onConnected    util.Subject[UserInfo]
	onDisconnected util.Subject[bool]
}

var networkClientInstance networkClient

func NetworkClient() *networkClient {
	return &networkClientInstance
}

func InitNetworkClient() error {
	const kStorage = "storage"
	networkClientInstance = networkClient{
		configDir:      filepath.Join(kStorage, "config.json"),
		config:         networkClientConfig{ServerURL: "http://132.145.98.4/"},
		onConnected:    util.MakeSubject[UserInfo](),
		onDisconnected: util.MakeSubject[bool](),
	}

	data, err := os.ReadFile(networkClientInstance.configDir)

	//create default config
	if errors.Is(err, os.ErrNotExist) {
		return networkClientInstance.save()
	}

	//unexpected error
	if err != nil {
		return err
	}

	//read config
	return json.Unmarshal(data, &networkClientInstance.config)
}

func (c *networkClient) OnConnected() util.Subject[UserInfo] {
	return c.onConnected
}

func (c *networkClient) OnDisconnected() util.Subject[bool] {
	return c.onDisconnected
}

func (c *networkClient) save() error {
	data, err := json.Marshal(c.config)
	if err != nil {
		return err
	}
	return os.WriteFile(c.configDir, data, 0600)
}

func (c *networkClient) setConfig(username string, password string) error {
	c.config.Username = username
	c.config.Password = password
	return c.save()
}

func (c *networkClient) sendPost(method string, endpoint string, username string, password string, content io.Reader) (*http.Response, error) {
	//populate URL and request fields
	url, err := url.JoinPath(c.config.ServerURL, endpoint)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url, content)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(username, password)

	//send request
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	//check status and return response
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

func (c *networkClient) Login(username string, password string) error {
	if _, err := c.sendPost(http.MethodGet, "login", username, password, nil); err != nil {
		UIClient().setStorage(newLocalStorage())
		c.onDisconnected.NotifyAll(true)
		return err
	}
	UIClient().setStorage(newRemoteStorage())
	c.onConnected.NotifyAll(UserInfo{Username: username})
	return c.setConfig(username, password)
}

func (c *networkClient) LoginWithConfig() {
	c.Login(c.config.Username, c.config.Password)
}

func (c *networkClient) Logout() {
	UIClient().setStorage(newLocalStorage())
	c.onDisconnected.NotifyAll(true)
}

func (c *networkClient) Register(username string, password string) error {
	if _, err := c.sendPost(http.MethodPost, "register", username, password, nil); err != nil {
		return err
	}
	return c.setConfig(username, password)
}

func (c *networkClient) getAllAlbums() ([]Album, error) {
	resp, err := c.sendPost(http.MethodGet, "downloadAll", c.config.Username, c.config.Password, nil)
	if err != nil {
		fyne.LogError("failed to download all albums", err)
		return nil, err
	}
	defer resp.Body.Close()

	var albums []Album
	if err := json.NewDecoder(resp.Body).Decode(&albums); err != nil {
		fyne.LogError("failed to decode albums", err)
		return nil, err
	}
	return albums, nil
}

func (c *networkClient) uploadAlbum(album Album) error {
	content := bytes.Buffer{}
	if err := json.NewEncoder(&content).Encode(album); err != nil {
		return err
	}

	_, err := c.sendPost(http.MethodPost, "upload", c.config.Username, c.config.Password, &content)
	if err != nil {
		fyne.LogError("failed to upload album", err)
		return err
	}
	return nil
}

func (c *networkClient) removeAlbum(key AlbumKey) error {
	_, err := c.sendPost(http.MethodPost, "remove", c.config.Username, c.config.Password, bytes.NewBufferString(key.String()))
	if err != nil {
		fyne.LogError("failed to remove album", err)
		return err
	}
	return nil
}

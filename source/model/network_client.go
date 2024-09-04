package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"meowyplayer/util"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	kAlbumKeyParam = "albumKey"
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
	sync.Mutex
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

func (c *networkClient) getURL(endpoint string) string {
	url, err := url.JoinPath(c.config.ServerURL, endpoint)
	if err != nil {
		log.Panicln(err)
	}
	return url
}

func (c *networkClient) sendRequest(method string, endpoint string, contentType string, content io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, c.getURL(endpoint), content)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.SetBasicAuth(c.config.Username, c.config.Password)

	//execute the request
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

func (c *networkClient) sendPostJson(endpoint string, object json.Marshaler) (*http.Response, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(object); err != nil {
		return nil, err
	}
	return c.sendRequest(http.MethodPost, endpoint, "application/json", &buf)
}

func (c *networkClient) sendPostForm(endpoint string, form url.Values) (*http.Response, error) {
	return c.sendRequest(http.MethodPost, endpoint, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
}

func (c *networkClient) sendPostEmpty(endpoint string) (*http.Response, error) {
	return c.sendRequest(http.MethodPost, endpoint, "", nil)
}

func (c *networkClient) Login(username string, password string) error {
	c.Lock()
	defer c.Unlock()
	oldUsername, oldPassword := c.config.Username, c.config.Password
	c.config.Username, c.config.Password = username, password

	if _, err := c.sendPostEmpty("login"); err != nil {
		//restore the old config if fails to login
		c.config.Username, c.config.Password = oldUsername, oldPassword
		return err
	}

	if err := c.setConfig(username, password); err != nil {
		return err
	}
	c.onConnected.NotifyAll(UserInfo{Username: username})
	return StorageClient().setStorage(newRemoteStorage())
}

func (c *networkClient) LoginWithConfig() {
	if err := c.Login(c.config.Username, c.config.Password); err != nil {
		c.Logout()
	}
}

func (c *networkClient) Logout() error {
	c.Lock()
	defer c.Unlock()
	if err := c.setConfig("", ""); err != nil {
		return err
	}
	c.onDisconnected.NotifyAll(true)
	return StorageClient().setStorage(newLocalStorage())
}

func (c *networkClient) Register(username string, password string) error {
	c.Lock()
	defer c.Unlock()
	_, err := c.sendPostEmpty("register")
	return err
}

func (c *networkClient) MigrateLocalToRemote() error {
	c.Lock()
	defer c.Unlock()
	albums, err := newLocalStorage().getAllAlbums()
	if err != nil {
		return err
	}

	for _, album := range albums {
		if err := StorageClient().UploadAlbum(album); err != nil {
			return err
		}
	}
	return nil
}

func (c *networkClient) getAllAlbums() ([]Album, error) {
	resp, err := c.sendPostEmpty("downloadAll")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var albums []Album
	err = json.NewDecoder(resp.Body).Decode(&albums)
	return albums, err
}

func (c *networkClient) uploadAlbum(album Album) error {
	_, err := c.sendPostJson("upload", album)
	return err
}

func (c *networkClient) removeAlbum(key AlbumKey) error {
	_, err := c.sendPostForm("remove", url.Values{kAlbumKeyParam: {key.String()}})
	return err
}

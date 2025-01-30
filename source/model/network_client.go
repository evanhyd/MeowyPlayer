package model

import (
	"bytes"
	"encoding/base64"
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
)

const (
	kAlbumKeyParam = "albumKey"
)

type UserProfile struct {
	Username string
}

type connectionConfig struct {
	Endpoint       string `json:"endpoint"`
	Username       string `json:"username"`
	PasswordBase64 string `json:"password"`
	// Token          []byte `json:"token"`
}

func (c *connectionConfig) getURL(resource string) string {
	URL, err := url.JoinPath(c.Endpoint, resource)
	if err != nil {
		log.Panicln(err)
	}
	return URL
}

func (c *connectionConfig) setAuth(username string, password string) {
	c.Username = username
	c.PasswordBase64 = base64.StdEncoding.EncodeToString([]byte(password))
}

func (c *connectionConfig) getAuth() (string, string) {
	password, _ := base64.StdEncoding.DecodeString(c.PasswordBase64)
	return c.Username, string(password)
}

type networkClient struct {
	configDir      string
	config         connectionConfig
	onConnected    util.Subject[UserProfile]
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
		config:         connectionConfig{Endpoint: "http://132.145.98.4/"},
		onConnected:    util.MakeSubject[UserProfile](),
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

func (c *networkClient) OnConnected() util.Subject[UserProfile] {
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

func (c *networkClient) sendRequest(method string, resource string, contentType string, content io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, c.config.getURL(resource), content)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.SetBasicAuth(c.config.getAuth())

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

func (c *networkClient) sendPostJson(resource string, object json.Marshaler) (*http.Response, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(object); err != nil {
		return nil, err
	}
	return c.sendRequest(http.MethodPost, resource, "application/json", &buf)
}

func (c *networkClient) sendPostForm(resource string, form url.Values) (*http.Response, error) {
	return c.sendRequest(http.MethodPost, resource, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
}

func (c *networkClient) sendPostEmpty(resource string) (*http.Response, error) {
	return c.sendRequest(http.MethodPost, resource, "", nil)
}

func (c *networkClient) LoginManually(username string, password string) error {
	// Try to login with new config. Revert if fails.
	oldConfig := c.config
	c.config.setAuth(username, password)
	if _, err := c.sendPostEmpty("login"); err != nil {
		c.config = oldConfig
		return err
	}

	// Login successfully, save the login credential.
	if err := c.save(); err != nil {
		return err
	}
	c.onConnected.NotifyAll(UserProfile{Username: username})
	return StorageClient().setStorage(newRemoteStorage())
}

func (c *networkClient) LoginWithConfig() error {
	if _, err := c.sendPostEmpty("login"); err != nil {
		c.Logout()
		return nil
	}
	c.onConnected.NotifyAll(UserProfile{Username: c.config.Username})
	return StorageClient().setStorage(newRemoteStorage())
}

func (c *networkClient) Logout() error {
	c.config.setAuth("", "")
	if err := c.save(); err != nil {
		return err
	}
	c.onDisconnected.NotifyAll(true)
	return StorageClient().setStorage(newLocalStorage())
}

func (c *networkClient) Register(username string, password string) error {
	// Try register with new username and password. Revert if fails.
	oldConfig := c.config
	c.config.setAuth(username, password)
	_, err := c.sendPostEmpty("register")
	if err != nil {
		c.config = oldConfig
	}
	return err
}

func (c *networkClient) UploadLocalToTheAccount() error {
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

func (c *networkClient) BackupAlbumsToLocal() error {
	albums, err := c.getAllAlbums()
	if err != nil {
		return err
	}

	localStorage := newLocalStorage()
	for _, album := range albums {
		if err := localStorage.uploadAlbum(album); err != nil {
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

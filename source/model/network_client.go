package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"meowyplayer/util"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type networkClientConfig struct {
	ServerURL string `json:"serverURL"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	// Token          []byte `json:"token"`
}

type networkClient struct {
	configDir           string
	config              networkClientConfig
	onConnectionChanged util.Subject[bool]
}

var networkClientInstance networkClient

func NetworkClient() *networkClient {
	return &networkClientInstance
}

func InitNetworkClient() error {
	const kStorage = "storage"
	networkClientInstance = networkClient{
		configDir:           filepath.Join(kStorage, "config.json"),
		config:              networkClientConfig{ServerURL: "http://132.145.98.4/"},
		onConnectionChanged: util.MakeSubject[bool](),
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

func (c *networkClient) OnConnectionChanged() util.Subject[bool] {
	return c.onConnectionChanged
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

func (c *networkClient) sendPost(method string, endpoint string, username string, password string) (*http.Response, error) {
	//populate URL and request fields
	url, err := url.JoinPath(c.config.ServerURL, endpoint)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url, nil)
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
	if _, err := c.sendPost(http.MethodPost, "login", username, password); err != nil {
		UIClient().setStorage(newLocalStorage())
		c.onConnectionChanged.NotifyAll(false)
		return err
	}
	UIClient().setStorage(newRemoteStorage())
	c.onConnectionChanged.NotifyAll(true)
	return c.setConfig(username, password)
}

func (c *networkClient) LoginWithConfig() {
	c.Login(c.config.Username, c.config.Password)
}

func (c *networkClient) Logout() {
	UIClient().setStorage(newLocalStorage())
	c.onConnectionChanged.NotifyAll(false)
}

func (c *networkClient) Register(username string, password string) error {
	if _, err := c.sendPost(http.MethodPost, "register", username, password); err != nil {
		return err
	}
	return c.setConfig(username, password)
}

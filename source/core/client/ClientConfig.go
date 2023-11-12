package client

var config = clientConfig{ServerUrl: "http://localhost", AutoBackup: false}

func Config() *clientConfig {
	return &config
}

type clientConfig struct {
	ServerUrl  string `json:"serverUrl"`
	AutoBackup bool   `json:"autoBackup"`
}

func (c *clientConfig) SetServer(url string) {
	c.ServerUrl = url
}

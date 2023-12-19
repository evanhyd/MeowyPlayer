package client

import (
	"errors"
	"io/fs"
	"os"
	"sync"

	"meowyplayer.com/core/resource"
	"meowyplayer.com/utility/pattern"
	"meowyplayer.com/utility/ujson"
)

var config = configManager{}

func Config() *configManager {
	return &config
}

type configManager struct {
	accessLock sync.Mutex
	config     pattern.Data[resource.Config]
}

func (c *configManager) Initialize() error {
	_, err := os.Stat(resource.ConfigFile())
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	//create default config
	if errors.Is(err, fs.ErrNotExist) {
		config := resource.Config{Name: "Guest", ServerUrl: "http://localhost"}
		if err := ujson.WriteFile(resource.ConfigFile(), config); err != nil {
			return err
		}
	}
	return c.load()
}

func (c *configManager) save() error {
	return ujson.WriteFile(resource.ConfigFile(), c.config.Get())
}

func (c *configManager) load() error {
	config := resource.Config{}
	if err := ujson.ReadFile(resource.ConfigFile(), &config); err != nil {
		return err
	}
	c.config.Set(config)
	return nil
}

func (c *configManager) Name() string {
	return c.config.Get().Name
}

func (c *configManager) ServerUrl() string {
	return c.config.Get().ServerUrl
}

func (c *configManager) SetName(name string) {
	c.accessLock.Lock()
	defer c.accessLock.Unlock()
	config := c.config.Get()
	config.Name = name
	c.config.Set(config)
	c.save()
}

func (c *configManager) SetServerUrl(url string) {
	c.accessLock.Lock()
	defer c.accessLock.Unlock()
	config := c.config.Get()
	config.ServerUrl = url
	c.config.Set(config)
	c.save()
}

func (c *configManager) AddListener(observer pattern.Observer[resource.Config]) {
	c.config.Attach(observer)
}

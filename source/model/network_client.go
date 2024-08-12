package model

type networkClient struct {
	storage *remoteStorage
}

var networkClientInstance networkClient

func NetworkClient() *networkClient {
	return &networkClientInstance
}

func InitNetworkClient() error {
	networkClientInstance = networkClient{storage: newRemoteStorage()}
	return networkClientInstance.storage.initialize()
}

func (c *networkClient) Run() error {
	//try log in
	return nil
}

func (c *networkClient) Login(username string, password string) error {
	if err := c.storage.authenticate(username, password); err != nil {
		return err
	}
	return c.storage.setConfig(username, password)
}

func (c *networkClient) Register(username string, password string) error {
	if err := c.storage.register(username, password); err != nil {
		return err
	}
	return c.storage.setConfig(username, password)
}

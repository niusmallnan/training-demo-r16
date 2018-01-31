package rancherapi

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/niusmallnan/training-demo/config"
	"github.com/rancher/go-rancher/v2"
)

type MyAPIClient struct {
	rclient *client.RancherClient
}

func NewClient(conf config.Config) (*MyAPIClient, error) {
	apiClient, err := client.NewRancherClient(&client.ClientOpts{
		Timeout:   time.Second * 30,
		Url:       conf.CattleURL,
		AccessKey: conf.CattleAccessKey,
		SecretKey: conf.CattleSecretKey,
	})
	if err != nil {
		return nil, err
	}

	return &MyAPIClient{apiClient}, nil
}

func (c *MyAPIClient) ListHost() error {
	listOpts := client.NewListOpts()
	collection, err := c.rclient.Host.List(listOpts)
	if err != nil {
		logrus.Errorf("Failed to request rancher-api: %v", err)
		return err
	}
	for _, h := range collection.Data {
		logrus.Infof("Got a host: %v", h)
	}

	return nil
}

package ranchermd

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/rancher/go-rancher-metadata/metadata"
)

const (
	changeCheckInterval = 5
	metadataURL         = "http://%s/2016-07-29"
)

type Watcher struct {
	m metadata.Client
}

func NewWatcher(metadataAddress string) (*Watcher, error) {
	m, err := metadata.NewClientAndWait(fmt.Sprintf(metadataURL, metadataAddress))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create metadata client")
	}

	o := &Watcher{m}

	return o, nil
}

func (o *Watcher) Start() {
	logrus.Infof("Rancher metadata watcher: Start")
	go o.m.OnChange(changeCheckInterval, o.onChangeNoError)
}

func (o *Watcher) onChangeNoError(version string) {
	if err := o.doYourJob(); err != nil {
		logrus.Errorf("Failed to apply your jobs: %v", err)
	}
}

func (o *Watcher) doYourJob() error {
	allContainers, err := o.m.GetContainers()
	if err != nil {
		return errors.Wrap(err, "Failed to get containers from metadata")
	}

	for _, c := range allContainers {
		logrus.Infof("Got a container: %s", c.Name)
	}

	return nil
}

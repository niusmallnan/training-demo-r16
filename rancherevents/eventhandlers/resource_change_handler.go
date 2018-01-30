package eventHandlers

import (
	log "github.com/Sirupsen/logrus"
	revents "github.com/rancher/event-subscriber/events"
	"github.com/rancher/go-rancher/client"
)

type ResourceChangeHandler struct {
}

func NewResourceChangeHandler() *ResourceChangeHandler {
	return &ResourceChangeHandler{}
}

func (h *ResourceChangeHandler) Handler(event *revents.Event, cli *client.RancherClient) error {
	if event.ResourceType == "service" {
		resource, _ := event.Data["resource"].(map[string]interface{})
		state, _ := resource["state"].(string)
		name, _ := resource["name"].(string)
		log.Infof("%s %s", name, state)

	}
	return nil
}

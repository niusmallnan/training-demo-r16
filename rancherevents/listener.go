package rancherevents

import (
	revents "github.com/rancher/event-subscriber/events"

	"github.com/niusmallnan/training-demo/config"
	reventhandlers "github.com/niusmallnan/training-demo/rancherevents/eventhandlers"
)

func ConnectToEventStream(conf config.Config) error {
	ehs := map[string]revents.EventHandler{
		"resource.change": reventhandlers.NewResourceChangeHandler().Handler,
	}
	router, err := revents.NewEventRouter("", 0, conf.CattleURL, conf.CattleAccessKey, conf.CattleSecretKey, nil, ehs, "", conf.WorkerCount, revents.DefaultPingConfig)
	if err != nil {
		return err
	}
	err = router.StartWithoutCreate(nil)
	return err
}

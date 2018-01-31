package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/niusmallnan/training-demo/config"
	"github.com/niusmallnan/training-demo/healthcheck"
	"github.com/niusmallnan/training-demo/rancherapi"
	"github.com/niusmallnan/training-demo/rancherevents"
	"github.com/niusmallnan/training-demo/ranchermd"
	"github.com/urfave/cli"
)

var VERSION = "v0.1.0-dev"

func main() {
	app := cli.NewApp()
	app.Name = "training-demo"
	app.Version = VERSION
	app.Usage = "You need help!"
	app.Action = launch

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug, d",
			EnvVar: "RANCHER_DEBUG",
		},
		cli.StringFlag{
			Name:   "cattle-url",
			Usage:  "URL for cattle API",
			EnvVar: "CATTLE_URL",
		},
		cli.StringFlag{
			Name:   "cattle-access-key",
			Usage:  "Cattle API Access Key",
			EnvVar: "CATTLE_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "cattle-secret-key",
			Usage:  "Cattle API Secret Key",
			EnvVar: "CATTLE_SECRET_KEY",
		},
		cli.IntFlag{
			Name:   "health-check-port",
			Value:  20220,
			Usage:  "Port to configure an HTTP health check listener on",
			EnvVar: "HEALTH_CHECK_PORT",
		},
		cli.IntFlag{
			Name:   "worker-count",
			Value:  50,
			Usage:  "Number of workers for handling events",
			EnvVar: "WORKER_COUNT",
		},
		cli.StringFlag{
			Name:   "metadata-address",
			Value:  config.DefaultMetadataAddress,
			EnvVar: "RANCHER_METADATA_ADDRESS",
		},
		cli.BoolFlag{
			Name:   "enable-metadata-watcher",
			EnvVar: "RANCHER_ENABLE_MD_WATCHER",
		},
	}
	app.Run(os.Args)
}

func launch(c *cli.Context) error {
	if c.Bool("debug") {
		os.Setenv("RANCHER_CLIENT_DEBUG", "true")
		log.SetLevel(log.DebugLevel)
	}
	conf := config.Conf(c)

	if c.Bool("enable-metadata-watcher") {
		watcher, err := ranchermd.NewWatcher(c.String("metadata-address"))
		if err != nil {
			log.Errorf("Failed to create rancher metadata watcher: %v", err)
		}
		watcher.Start()
	}

	api, err := rancherapi.NewClient(conf)
	if err != nil {
		log.Errorf("Failed to create rancher api client: %v", err)
	}
	api.ListHost()

	resultChan := make(chan error)

	go func(rc chan error) {
		err = rancherevents.ConnectToEventStream(conf)
		log.Errorf("Rancher stream listener exited with error: %v", err)
		rc <- err
	}(resultChan)

	go func(rc chan error) {
		err = healthcheck.StartHealthCheck(conf.HealthCheckPort)
		log.Errorf("HealthCheck exit with error: %v", err)
		rc <- err
	}(resultChan)

	err = <-resultChan
	log.Info("Exiting...")
	return err
}

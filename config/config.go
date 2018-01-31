package config

import "github.com/urfave/cli"

const (
	DefaultMetadataAddress = "169.254.169.250"
)

type Config struct {
	CattleURL       string
	CattleAccessKey string
	CattleSecretKey string
	WorkerCount     int
	HealthCheckPort int
}

func Conf(context *cli.Context) Config {
	config := Config{
		CattleURL:       context.String("cattle-url"),
		CattleAccessKey: context.String("cattle-access-key"),
		CattleSecretKey: context.String("cattle-secret-key"),
		WorkerCount:     context.Int("worker-count"),
		HealthCheckPort: context.Int("health-check-port"),
	}
	return config
}

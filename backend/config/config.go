package config

import (
	"github.com/kelseyhightower/envconfig"
)

var DefaultConfig = &Config{
	Database: Database{
		Driver:   DatabaseDriverPostgres,
		Host:     "localhost",
		Port:     5432,
		User:     "s3ase",
		Password: "s3ase",
		Database: "s3ase",
		Options:  "sslmode=disable&connect_timeout=30",
	},
	Redis: Redis{
		Host:               "localhost",
		Port:               6379,
		Username:           "",
		Password:           "",
		MinIdleConnections: 0,
		MaxRetries:         3,
	},
	Aws: Aws{
		DefaultRegion: "eu-west-1",
	},
	Server: Server{
		Port: 5555,
		SSL:  false,
	},
	Auth: Auth{
		RedirectUrl:   "http://localhost:5555",
		UIRedirectUrl: "http://localhost:3000/auth/callback",
	},
	Email: Email{
		FromAddress: "S3ase <no-reply@s3ase.dev>",
	},
	Job: Job{
		Concurrency: 10,
	},
	Pubsub: Pubsub{
		App:            "s3ase",
		Namespace:      "s3ase",
		HealthInterval: 2,
		SendTimeout:    60,
		ChannelSize:    100,
	},
}

func LoadConfig() (*Config, error) {
	config := DefaultConfig

	// Override config from environment variables
	err := envconfig.Process("", config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

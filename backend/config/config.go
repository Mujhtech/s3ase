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
		Host:     "localhost",
		Port:     6379,
		Username: "",
		Password: "",
	},
	Aws: Aws{
		Region: "eu-west-1",
	},
	Server: Server{
		Port: 5555,
		SSL:  false,
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

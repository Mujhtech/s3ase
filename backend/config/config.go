package config

import (
	"fmt"
	"regexp"
	"time"

	"github.com/kelseyhightower/envconfig"
)

var r2AccountIDPattern = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

var DefaultConfig = &Config{
	Environment:       "development",
	EncryptionKey:     "development-only-change-me-key!!",
	DomainCnameTarget: "files.s3ase.dev",
	Cache: Cache{
		Provider: CacheProviderRedis,
	},
	ObjectStorage: ObjectStorage{
		Provider: ObjectStorageProviderS3,
	},
	Cors: Cors{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Accept", "Authorization", "Accept-Encoding", "Content-Length", "Content-Type", "X-CSRF-Token", "X-Requested-With", "X-Requested-Id", "x-app-id", "Tus-Resumable", "Upload-Length", "Upload-Offset", "Upload-Metadata", "Upload-Defer-Length", "Upload-Concat"},
		ExposedHeaders:   []string{"Link", "Location", "Tus-Resumable", "Tus-Version", "Tus-Extension", "Tus-Max-Size", "Upload-Length", "Upload-Offset", "Upload-Metadata", "X-File-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	},
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
		DB:                 1,
	},
	Aws: Aws{
		DefaultRegion: "eu-west-2",
		UsePathStyle:  false,
	},
	Server: Server{
		Port:    5555,
		SSL:     false,
		Timeout: 30,
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
		Provider:       PubsubProviderInMemory,
		App:            "s3ase",
		Namespace:      "s3ase",
		HealthInterval: 2,
		SendTimeout:    60,
		ChannelSize:    500,
	},
	Protocol: Protocol{
		MaxSize:                1073741824,
		UploadProgressInterval: time.Second,
		NetworkTimeout:         30 * time.Second,
	},
}

func LoadConfig() (*Config, error) {
	config := *DefaultConfig

	// Override config from environment variables
	err := envconfig.Process("", &config)
	if err != nil {
		return nil, err
	}

	if err = config.validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) validate() error {
	if len(c.EncryptionKey) != 32 {
		return fmt.Errorf("encryption key must be exactly 32 bytes")
	}
	if c.Environment == "production" && c.EncryptionKey == DefaultConfig.EncryptionKey {
		return fmt.Errorf("encryption key must be changed in production")
	}
	if c.Server.Port == 0 {
		return fmt.Errorf("server port cannot be zero")
	}
	if c.DomainCnameTarget == "" {
		return fmt.Errorf("domain CNAME target cannot be empty")
	}
	if c.Protocol.MaxSize <= 0 {
		return fmt.Errorf("protocol max size must be greater than zero")
	}
	if c.Protocol.UploadProgressInterval <= 0 {
		return fmt.Errorf("protocol upload progress interval must be greater than zero")
	}
	if c.Protocol.NetworkTimeout <= 0 {
		return fmt.Errorf("protocol network timeout must be greater than zero")
	}
	if (c.Auth.GithubAuth.ClientID == "") != (c.Auth.GithubAuth.ClientSecret == "") {
		return fmt.Errorf("github auth client id and secret must be configured together")
	}
	if (c.Auth.GoogleAuth.ClientID == "") != (c.Auth.GoogleAuth.ClientSecret == "") {
		return fmt.Errorf("google auth client id and secret must be configured together")
	}
	switch c.ObjectStorage.Provider {
	case ObjectStorageProviderS3:
		if (c.Aws.AccessKey == "") != (c.Aws.SecretKey == "") {
			return fmt.Errorf("aws access key and secret key must be configured together")
		}
	case ObjectStorageProviderR2:
		if !r2AccountIDPattern.MatchString(c.R2.AccountID) {
			return fmt.Errorf("r2 account id must contain only letters and numbers")
		}
		if c.R2.AccessKeyID == "" || c.R2.SecretKey == "" {
			return fmt.Errorf("r2 access key id and secret access key are required")
		}
		if c.R2.Jurisdiction != "" && c.R2.Jurisdiction != "eu" && c.R2.Jurisdiction != "fedramp" {
			return fmt.Errorf("r2 jurisdiction must be empty, eu, or fedramp")
		}
	default:
		return fmt.Errorf("object storage provider must be s3 or r2")
	}

	// Validate database configuration
	if c.Database.Host == "" {
		return fmt.Errorf("database host cannot be empty")
	}
	if c.Database.Port == 0 {
		return fmt.Errorf("database port cannot be zero")
	}

	dbDsn := c.Database.BuildDsn()
	if dbDsn == "" {
		return fmt.Errorf("database dsn is empty")
	}

	return nil
}

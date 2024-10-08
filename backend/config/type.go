package config

import (
	"fmt"
	"net/url"
)

type DatabaseDriver string

const (
	DatabaseDriverPostgres DatabaseDriver = "postgres"
	DatabaseDriverSqlite3  DatabaseDriver = "sqlite3"

	DefaultConfigFilePath string = ".env"
)

type Config struct {
	Database Database `json:"database"`
	Redis    Redis    `json:"redis"`
	Aws      Aws      `json:"aws"`
	Server   Server   `json:"server"`
}

// Database defines database configuration
type Database struct {
	Driver   DatabaseDriver `json:"driver" envconfig:"DB_DRIVER"`
	Host     string         `json:"host" envconfig:"DB_HOST"`
	Port     int            `json:"port" envconfig:"DB_PORT"`
	User     string         `json:"user" envconfig:"DB_USER"`
	Password string         `json:"password" envconfig:"DB_PASSWORD"`
	Database string         `json:"database" envconfig:"DB_DATABASE"`
	Options  string         `json:"options" envconfig:"DB_OPTIONS"`
}

// Aws defines AWS configuration
type Aws struct {
	Region string `envconfig:"AWS_REGION"`
}

// Server defines server configuration
type Server struct {
	Port        uint32 `json:"port" envconfig:"PORT"`
	SSL         bool   `json:"ssl" envconfig:"SSL"`
	SSLCertFile string `json:"ssl_cert_file" envconfig:"SSL_CERT_FILE"`
	SSLKeyFile  string `json:"ssl_key_file" envconfig:"SSL_KEY_FILE"`
	Timeout     uint32 `json:"timeout" envconfig:"TIMEOUT"`
}

// Redis defines redis configuration
type Redis struct {
	Host     string `json:"host" envconfig:"REDIS_HOST"`
	Port     int    `json:"port" envconfig:"REDIS_PORT"`
	Username string `json:"username" envconfig:"REDIS_USERNAME"`
	Password string `json:"password" envconfig:"REDIS_PASSWORD"`
}

func (d *Database) BuildDsn() string {
	if d.Driver == "" {
		return ""
	}

	authPart := ""
	if d.User != "" || d.Password != "" {
		authPrefix := url.UserPassword(d.User, d.Password)
		authPart = fmt.Sprintf("%s@", authPrefix)
	}

	dbPart := ""
	if d.Database != "" {
		dbPart = fmt.Sprintf("/%s", d.Database)
	}

	optPart := ""
	if d.Options != "" {
		optPart = fmt.Sprintf("?%s", d.Options)
	}

	return fmt.Sprintf("%s://%s%s:%d%s%s", d.Driver, authPart, d.Host, d.Port, dbPart, optPart)
}

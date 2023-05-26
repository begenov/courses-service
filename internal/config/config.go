package config

import (
	"fmt"
	"os"
	"time"

	"github.com/subosito/gotenv"
)

const (
	defaultServerPort               = "8080"
	defaultServerRWTimeout          = 10 * time.Second
	defaultServerMaxHeaderMegabytes = 1
	defaultAccessTokenTTL           = 15 * time.Minute
	defaultRefreshTokenTTL          = 24 * time.Hour * 30
)

type Config struct {
	Server   serverConfig
	Database databaseConfig
}

type serverConfig struct {
	Port               string
	ReadTImeout        time.Duration
	WriteTimeout       time.Duration
	MaxHeaderMegabytes int
}

type databaseConfig struct {
	Driver string
	DSN    string
}

func Init(path string) (*Config, error) {
	err := gotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load environment variables from file: %v", err)
	}
	driver := os.Getenv("DRIVER")
	dsn := os.Getenv("DSN_COURSES")
	return &Config{
		Server: serverConfig{
			Port:               defaultServerPort,
			ReadTImeout:        defaultServerRWTimeout,
			WriteTimeout:       defaultServerRWTimeout,
			MaxHeaderMegabytes: defaultServerMaxHeaderMegabytes,
		},
		Database: databaseConfig{
			DSN:    dsn,
			Driver: driver,
		},
	}, nil
}

package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/subosito/gotenv"
)

const (
	defaultServerPort               = "8080"
	defaultServerRWTimeout          = 10 * time.Second
	defaultServerMaxHeaderMegabytes = 1
	TTLCache                        = 24 * time.Hour * 30
)

type Config struct {
	Server   serverConfig
	Database databaseConfig
	Redis    RedisConfig
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

type RedisConfig struct {
	DB       int
	Host     string
	Port     string
	Password string
	Ttl      time.Duration
}

func Init(path string) (*Config, error) {
	err := gotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load environment variables from file: %v", err)
	}
	driver := os.Getenv("DRIVER")
	dsn := os.Getenv("DSN_COURSES")

	host := os.Getenv("HOST_REDIS")
	port := os.Getenv("PORT_REDIS")
	password_redis := os.Getenv("PASSWORD_REDIS")
	db_redis, err := strconv.Atoi(os.Getenv("DB_REDIS"))
	if err != nil {
		return nil, err
	}
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
		Redis: RedisConfig{
			Host:     host,
			DB:       db_redis,
			Password: password_redis,
			Port:     port,
			Ttl:      TTLCache,
		},
	}, nil
}

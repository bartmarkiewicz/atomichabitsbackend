package config

import (
	"log"
	"time"

	"github.com/joeshaw/envdecode"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}
type ServerConfig struct {
	Port         int           `env:"SERVER_PORT,required"`
	TimeoutRead  time.Duration `env:"SERVER_TIMEOUT_READ,required"`
	TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE,required"`
	TimeoutIdle  time.Duration `env:"SERVER_TIMEOUT_IDLE,required"`
	Debug        bool          `env:"SERVER_DEBUG,required"`
}

type DatabaseConfig struct {
	Host         string `env:"DB_HOST,required"`
	Port         int    `env:"DB_PORT,required"`
	Username     string `env:"DB_USER,required"`
	Password     string `env:"DB_PASS,required"`
	DatabaseName string `env:"DB_NAME,required"`
	Debug        bool   `env:"DB_DEBUG,required"`
}

func New() *Config {
	var c Config
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	return &c
}

func NewDB() *DatabaseConfig {
	var databaseConfig DatabaseConfig
	if err := envdecode.StrictDecode(&databaseConfig); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	return &databaseConfig
}

package config

import (
	"flag"
	"fmt"
	"os"
	"sync"
)

const (
	defaultServerAddress = ":8080"
	defaultBaseURL       = "http://localhost:8080"
	defaultDatabaseDSN   = ""
)

type Config struct {
	ServerAddress string
	BaseURL       string
	DatabaseDSN   string
	Restore       bool
	SyncWrite     bool
}

var (
	instance *Config
	once     sync.Once
)

// New returns a singleton instance of Config. It parses command line flags and environment variables
// only once, even if called multiple times.
func New() (*Config, error) {
	once.Do(func() {
		cfg := Config{}

		flag.BoolVar(&cfg.Restore, "r", false, "restore from file")
		flag.StringVar(&cfg.DatabaseDSN, "d", defaultDatabaseDSN, "database DSN")
		flag.StringVar(&cfg.ServerAddress, "a", defaultServerAddress, "server address")
		flag.StringVar(&cfg.BaseURL, "b", defaultBaseURL, "base URL")
		flag.BoolVar(&cfg.SyncWrite, "s", false, "sync write")
		flag.Parse()

		if os.Getenv("SERVER_ADDRESS") != "" {
			cfg.ServerAddress = os.Getenv("SERVER_ADDRESS")
		}
		if os.Getenv("BASE_URL") != "" {
			cfg.BaseURL = os.Getenv("BASE_URL")
		}

		// Use environment variables for database connection if available
		if os.Getenv("DB_HOST") != "" {
			cfg.DatabaseDSN = fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
				os.Getenv("DB_USER"),
				os.Getenv("DB_PASSWORD"),
				os.Getenv("DB_HOST"),
				os.Getenv("DB_PORT"),
				os.Getenv("DB_NAME"),
			)
		} else if os.Getenv("DATABASE_DSN") != "" {
			cfg.DatabaseDSN = os.Getenv("DATABASE_DSN")
		}
		instance = &cfg
	})
	return instance, nil
}

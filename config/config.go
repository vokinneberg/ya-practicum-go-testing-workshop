package config

import (
	"flag"
	"os"
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

// New parses command line flags and environment variables and returns a new Config. If no en
func New() (*Config, error) {
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
	if os.Getenv("DATABASE_DSN") != "" {
		cfg.DatabaseDSN = os.Getenv("DATABASE_DSN")
	}
	if os.Getenv("RESTORE") != "" {
		cfg.Restore = os.Getenv("RESTORE") == "true"
	}
	if os.Getenv("SYNC_WRITE") != "" {
		cfg.SyncWrite = os.Getenv("SYNC_WRITE") == "true"
	}
	return &cfg, nil
}

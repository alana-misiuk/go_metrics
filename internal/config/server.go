package config

import (
	"flag"
	"os"
	"strconv"
)

type ServerConfig struct {
	Addr            string
	FlagLogLevel    string
	StoreInterval   int
	FileStoragePath string
	Restore         bool
}

func ParseServer(args []string) (ServerConfig, error) {
	var cfg ServerConfig

	fs := flag.NewFlagSet("server", flag.ContinueOnError)
	fs.StringVar(&cfg.Addr, "a", "localhost:8080", "server address")
	fs.StringVar(&cfg.FlagLogLevel, "l", "info", "log level")
	fs.IntVar(&cfg.StoreInterval, "i", 300, "store interval in seconds (0 = synchronous)")
	fs.StringVar(&cfg.FileStoragePath, "f", "/tmp/metrics-db.json", "file storage path")
	fs.BoolVar(&cfg.Restore, "r", true, "restore metrics on startup")

	if err := fs.Parse(args); err != nil {
		return ServerConfig{}, err
	}

	if v := os.Getenv("ADDRESS"); v != "" {
		cfg.Addr = v
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		cfg.FlagLogLevel = envLogLevel
	}
	if err := intFromEnv("STORE_INTERVAL", &cfg.StoreInterval); err != nil {
		return ServerConfig{}, err
	}
	if v := os.Getenv("FILE_STORAGE_PATH"); v != "" {
		cfg.FileStoragePath = v
	}
	if v := os.Getenv("RESTORE"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return ServerConfig{}, err
		}
		cfg.Restore = b
	}

	return cfg, nil
}

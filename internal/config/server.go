package config

import (
	"flag"
	"os"
)

type ServerConfig struct {
	Addr         string
	FlagLogLevel string
}

func ParseServer(args []string) (ServerConfig, error) {
	var cfg ServerConfig

	fs := flag.NewFlagSet("server", flag.ContinueOnError)
	fs.StringVar(&cfg.Addr, "a", "localhost:8080", "server address")
	fs.StringVar(&cfg.FlagLogLevel, "l", "info", "log level")

	if err := fs.Parse(args); err != nil {
		return ServerConfig{}, err
	}

	if v := os.Getenv("ADDRESS"); v != "" {
		cfg.Addr = v
	}

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		cfg.FlagLogLevel = envLogLevel
	}

	return cfg, nil
}

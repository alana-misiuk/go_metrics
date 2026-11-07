package config

import "flag"

type ServerConfig struct {
	Addr string
}

func ParseServer(args []string) (ServerConfig, error) {
	var cfg ServerConfig

	fs := flag.NewFlagSet("server", flag.ContinueOnError)
	fs.StringVar(&cfg.Addr, "a", "localhost:8080", "server address")

	if err := fs.Parse(args); err != nil {
		return ServerConfig{}, err
	}
	return cfg, nil
}

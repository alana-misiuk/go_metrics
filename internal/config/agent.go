package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type CollectorConfig struct {
	PollInterval int
}

type SenderConfig struct {
	Addr           string
	ReportInterval int
}

type AgentConfig struct {
	Collector CollectorConfig
	Sender    SenderConfig
}

func ParseAgent(args []string) (AgentConfig, error) {
	var cfg AgentConfig

	fs := flag.NewFlagSet("agent", flag.ContinueOnError)
	fs.StringVar(&cfg.Sender.Addr, "a", "localhost:8080", "server address")
	fs.IntVar(&cfg.Collector.PollInterval, "p", 2, "poll interval(s)")
	fs.IntVar(&cfg.Sender.ReportInterval, "r", 10, "report interval(s)")

	if err := fs.Parse(args); err != nil {
		return AgentConfig{}, err
	}

	if v := os.Getenv("ADDRESS"); v != "" {
		cfg.Sender.Addr = v
	}
	if err := intFromEnv("REPORT_INTERVAL", &cfg.Sender.ReportInterval); err != nil {
		return AgentConfig{}, err
	}
	if err := intFromEnv("POLL_INTERVAL", &cfg.Collector.PollInterval); err != nil {
		return AgentConfig{}, err
	}

	return cfg, nil
}

func intFromEnv(key string, dst *int) error {
	v := os.Getenv(key)
	if v == "" {
		return nil
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fmt.Errorf("%s: %w", key, err)
	}
	*dst = n
	return nil
}

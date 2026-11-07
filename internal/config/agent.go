package config

import "flag"

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
	return cfg, nil
}

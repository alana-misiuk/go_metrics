package main

import (
	"errors"
	"flag"
	"fmt"
	"go_metrics/internal/agent"
	"go_metrics/internal/config"
	"os"
	"time"
)

func main() {
	cfg, err := config.ParseAgent(os.Args[1:])
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			os.Exit(0)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	stats := &agent.AgentStats{}

	go func() {
		for {
			agent.CollectMetrics(stats)
			time.Sleep(time.Duration(cfg.Collector.PollInterval) * time.Second)
		}
	}()

	go func() {
		for {
			agent.SendMetrics(stats, cfg.Sender)
			time.Sleep(time.Duration(cfg.Sender.ReportInterval) * time.Second)
		}
	}()

	select {}
}

package main

import (
	"flag"
	"fmt"
	"go_metrics/internal/agent"
	"time"
)

func main() {
	addr := flag.String("a", "localhost:8080", "server address")
	pollInterval := flag.Int("p", 2, "poll interval in seconds")
	reportInterval := flag.Int("r", 10, "report interval in seconds")
	flag.Parse()

	serverURL := fmt.Sprintf("http://%s", *addr)
	stats := &agent.AgentStats{}

	go func() {
		for {
			agent.CollectMetrics(stats)
			time.Sleep(time.Duration(*pollInterval) * time.Second)
		}
	}()

	go func() {
		for {
			agent.SendMetrics(stats, serverURL)
			time.Sleep(time.Duration(*reportInterval) * time.Second)
		}
	}()

	select {}
}

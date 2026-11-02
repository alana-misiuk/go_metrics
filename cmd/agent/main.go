package main

import (
	"go_metrics/internal/agent"
	"time"
)

func main() {
    stats := &agent.AgentStats{}
    
    go func() {
		for {
			agent.CollectMetrics(stats)
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		for {
			agent.SendMetrics(stats)
			time.Sleep(10 * time.Second)
		}
	}()

	select {}

}

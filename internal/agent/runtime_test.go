package agent

import (
	"testing"
)

func TestCollectMetrics_PollCountIncrements(t *testing.T) {
	stats := &AgentStats{}
	CollectMetrics(stats)
	CollectMetrics(stats)

	if stats.PollCount != 2 {
		t.Errorf("expected PollCount=2, got %d", stats.PollCount)
	}
}

func TestCollectMetrics_AllocNonZero(t *testing.T) {
	stats := &AgentStats{}
	CollectMetrics(stats)

	if stats.Alloc == 0 {
		t.Error("expected Alloc > 0")
	}
}
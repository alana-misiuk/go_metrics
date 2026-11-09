package agent

import (
	"go_metrics/internal/config"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
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

type fakeServer struct {
	mu        sync.Mutex
	status    int
	received  []int64
	onCounter func()
}

func (f *fakeServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const prefix = "/update/counter/PollCount/"
	if strings.HasPrefix(r.URL.Path, prefix) {
		v, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, prefix), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		f.mu.Lock()
		f.received = append(f.received, v)
		f.mu.Unlock()
		if f.onCounter != nil {
			f.onCounter()
		}
	}
	w.WriteHeader(f.status)
}

func startFakeServer(t *testing.T, f *fakeServer) config.SenderConfig {
	t.Helper()
	srv := httptest.NewServer(f)
	t.Cleanup(srv.Close)
	return config.SenderConfig{Addr: strings.TrimPrefix(srv.URL, "http://")}
}

func TestSendMetrics_SendsDeltaAndResets(t *testing.T) {
	f := &fakeServer{status: http.StatusOK}
	cfg := startFakeServer(t, f)
	stats := &AgentStats{}

	CollectMetrics(stats)
	CollectMetrics(stats)
	SendMetrics(stats, cfg)

	CollectMetrics(stats)
	SendMetrics(stats, cfg)

	want := []int64{2, 1}
	if len(f.received) != len(want) || f.received[0] != want[0] || f.received[1] != want[1] {
		t.Errorf("server received %v, want %v", f.received, want)
	}
	if stats.PollCount != 0 {
		t.Errorf("PollCount after successful send = %d, want 0", stats.PollCount)
	}
}

func TestSendMetrics_KeepsCountOnServerError(t *testing.T) {
	f := &fakeServer{status: http.StatusInternalServerError}
	cfg := startFakeServer(t, f)
	stats := &AgentStats{}

	CollectMetrics(stats)
	CollectMetrics(stats)
	SendMetrics(stats, cfg)

	if stats.PollCount != 2 {
		t.Errorf("PollCount after failed send = %d, want 2", stats.PollCount)
	}
}

func TestSendMetrics_KeepsCountWhenServerDown(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	cfg := config.SenderConfig{Addr: strings.TrimPrefix(srv.URL, "http://")}
	srv.Close()
	stats := &AgentStats{}

	CollectMetrics(stats)
	SendMetrics(stats, cfg)

	if stats.PollCount != 1 {
		t.Errorf("PollCount when server is down = %d, want 1", stats.PollCount)
	}
}

func TestSendMetrics_KeepsPollsCollectedDuringSend(t *testing.T) {
	stats := &AgentStats{}
	f := &fakeServer{status: http.StatusOK}
	// Опрос случается, пока агент ждёт ответа на отправку счётчика.
	f.onCounter = func() { CollectMetrics(stats) }
	cfg := startFakeServer(t, f)

	CollectMetrics(stats)
	CollectMetrics(stats)
	SendMetrics(stats, cfg)

	if stats.PollCount != 1 {
		t.Errorf("PollCount = %d, want 1: poll made during send was lost", stats.PollCount)
	}
}

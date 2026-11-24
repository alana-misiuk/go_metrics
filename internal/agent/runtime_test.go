package agent

import (
	"compress/gzip"
	"encoding/json"
	"go_metrics/internal/config"
	models "go_metrics/internal/model"
	"io"
	"net/http"
	"net/http/httptest"
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
	metrics   []models.Metrics // все принятые метрики
	received  []int64          // приращения PollCount
	badReqs   []string         // запросы, не похожие на JSON-протокол
	onCounter func()
}

func (f *fakeServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var m models.Metrics

	body := r.Body
	if r.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			f.mu.Lock()
			f.badReqs = append(f.badReqs, "bad gzip")
			f.mu.Unlock()
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer gz.Close()
		body = io.NopCloser(gz)
	}

	if r.Method != http.MethodPost ||
		r.URL.Path != "/update" ||
		r.Header.Get("Content-Type") != "application/json" ||
		json.NewDecoder(body).Decode(&m) != nil {
		f.mu.Lock()
		f.badReqs = append(f.badReqs, r.Method+" "+r.URL.Path+" "+r.Header.Get("Content-Type"))
		f.mu.Unlock()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	isPollCount := m.ID == "PollCount" && m.MType == models.Counter && m.Delta != nil
	f.mu.Lock()
	f.metrics = append(f.metrics, m)
	if isPollCount {
		f.received = append(f.received, *m.Delta)
	}
	f.mu.Unlock()

	if isPollCount && f.onCounter != nil {
		f.onCounter()
	}
	w.WriteHeader(f.status)
}

func startFakeServer(t *testing.T, f *fakeServer) config.SenderConfig {
	t.Helper()
	srv := httptest.NewServer(f)
	t.Cleanup(srv.Close)
	return config.SenderConfig{Addr: strings.TrimPrefix(srv.URL, "http://")}
}

func TestSendMetrics_SendsJSON(t *testing.T) {
	f := &fakeServer{status: http.StatusOK}
	cfg := startFakeServer(t, f)
	stats := &AgentStats{}

	CollectMetrics(stats)
	SendMetrics(stats, cfg)

	if len(f.badReqs) > 0 {
		t.Fatalf("requests do not follow the JSON protocol: %v", f.badReqs)
	}

	const wantCount = 29
	if len(f.metrics) != wantCount {
		t.Errorf("sent %d metrics, want %d", len(f.metrics), wantCount)
	}

	sent := make(map[string]models.Metrics, len(f.metrics))
	for _, m := range f.metrics {
		sent[m.ID] = m
	}

	for _, name := range []string{"Alloc", "GCCPUFraction", "LastGC", "RandomValue", "TotalAlloc"} {
		m, ok := sent[name]
		if !ok {
			t.Errorf("gauge %s was not sent", name)
			continue
		}
		if m.MType != models.Gauge || m.Value == nil || m.Delta != nil {
			t.Errorf("gauge %s sent as %+v, want type gauge with value only", name, m)
		}
	}

	pc, ok := sent["PollCount"]
	if !ok {
		t.Fatal("PollCount was not sent")
	}
	if pc.MType != models.Counter || pc.Delta == nil || *pc.Delta != 1 || pc.Value != nil {
		t.Errorf("PollCount sent as %+v, want type counter with delta 1", pc)
	}
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

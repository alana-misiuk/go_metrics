package handler

import (
	"context"
	"go_metrics/internal/repository"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func newTestHandler() *Handler {
	return NewHandler(repository.NewMemStorage())
}

func withChiParams(r *http.Request, params map[string]string) *http.Request {
	rctx := chi.NewRouteContext()
	for k, v := range params {
		rctx.URLParams.Add(k, v)
	}
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

func TestUpdateMetric_GaugeOK(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodPost, "/update/gauge/Alloc/1.5", nil)
	req = withChiParams(req, map[string]string{"type": "gauge", "name": "Alloc", "value": "1.5"})
	w := httptest.NewRecorder()

	h.UpdateMetric(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestUpdateMetric_CounterOK(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodPost, "/update/counter/hits/10", nil)
	req = withChiParams(req, map[string]string{"type": "counter", "name": "hits", "value": "10"})
	w := httptest.NewRecorder()

	h.UpdateMetric(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestUpdateMetric_InvalidType(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodPost, "/update/unknown/Alloc/1.5", nil)
	req = withChiParams(req, map[string]string{"type": "unknown", "name": "Alloc", "value": "1.5"})
	w := httptest.NewRecorder()

	h.UpdateMetric(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUpdateMetric_InvalidValue(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodPost, "/update/gauge/Alloc/notanumber", nil)
	req = withChiParams(req, map[string]string{"type": "gauge", "name": "Alloc", "value": "notanumber"})
	w := httptest.NewRecorder()

	h.UpdateMetric(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUpdateMetric_WithoutName(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodPost, "/update/gauge//1.5", nil)
	req = withChiParams(req, map[string]string{"type": "gauge", "name": "", "value": "1.5"})
	w := httptest.NewRecorder()

	h.UpdateMetric(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestUpdateMetric_WrongMethod(t *testing.T) {
	r := NewRouter(newTestHandler())
	req := httptest.NewRequest(http.MethodGet, "/update/gauge/Alloc/1.5", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

func TestGetAllMetrics_ListsMetrics(t *testing.T) {
	storage := repository.NewMemStorage()
	storage.UpdateGauge("Alloc", 1.5)
	storage.UpdateCounter("PollCount", 3)
	h := NewHandler(storage)

	w := httptest.NewRecorder()
	h.GetAllMetrics(w, httptest.NewRequest(http.MethodGet, "/", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Errorf("unexpected Content-Type %q", ct)
	}
	body := w.Body.String()
	for _, want := range []string{"<li>Alloc (gauge) = 1.5</li>", "<li>PollCount (counter) = 3</li>"} {
		if !strings.Contains(body, want) {
			t.Errorf("body does not contain %q:\n%s", want, body)
		}
	}
}

func TestGetAllMetrics_EscapesNames(t *testing.T) {
	storage := repository.NewMemStorage()
	storage.UpdateGauge(`<script>alert("x")</script>&`, 1)
	h := NewHandler(storage)

	w := httptest.NewRecorder()
	h.GetAllMetrics(w, httptest.NewRequest(http.MethodGet, "/", nil))

	body := w.Body.String()
	if strings.Contains(body, "<script>") {
		t.Errorf("metric name was not escaped:\n%s", body)
	}
	if !strings.Contains(body, "&lt;script&gt;") || !strings.Contains(body, "&amp;") {
		t.Errorf("expected escaped name in body:\n%s", body)
	}
}

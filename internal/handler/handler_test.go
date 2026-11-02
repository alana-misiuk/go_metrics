package handler

import (
	"context"
	"go_metrics/internal/repository"
	"net/http"
	"net/http/httptest"
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
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/update/gauge/Alloc/1.5", nil)
	req = withChiParams(req, map[string]string{"type": "gauge", "name": "Alloc", "value": "1.5"})
	w := httptest.NewRecorder()

	h.UpdateMetric(w, req)

	if w.Code == http.StatusOK {
		t.Error("expected non-200 for GET request")
	}
}

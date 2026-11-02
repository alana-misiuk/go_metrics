package handler

import (
	"go_metrics/internal/repository"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestHandler() *Handler {
	return NewHandler(repository.NewMemStorage())
}

func TestUpdateMetric_GaugeOK(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodPost, "/update/gauge/Alloc/1.5", nil)
	w := httptest.NewRecorder()

	h.UpdateMetric(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestUpdateMetric_CounterOK(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodPost, "/update/counter/hits/10", nil)
	w := httptest.NewRecorder()

	h.UpdateMetric(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestUpdateMetric_InvalidType(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodPost, "/update/unknown/Alloc/1.5", nil)
	w := httptest.NewRecorder()

	h.UpdateMetric(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUpdateMetric_InvalidValue(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodPost, "/update/gauge/Alloc/notanumber", nil)
	w := httptest.NewRecorder()

	h.UpdateMetric(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUpdateMetric_WithoutName(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodPost, "/update/gauge/", nil)
	w := httptest.NewRecorder()

	h.UpdateMetric(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestUpdateMetric_WrongMethod(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/update/gauge/Alloc/1.5", nil)
	w := httptest.NewRecorder()

	h.UpdateMetric(w, req)

	if w.Code == http.StatusOK {
		t.Error("expected non-200 for GET request")
	}
}

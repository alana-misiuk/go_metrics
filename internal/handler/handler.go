package handler

import (
	"fmt"
	"go_metrics/internal/repository"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	storage repository.Storage
}

func NewHandler(storage repository.Storage) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")
	rawValue := chi.URLParam(r, "value")

	switch metricType {
	case "gauge":
		value, err := strconv.ParseFloat(rawValue, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h.storage.UpdateGauge(name, value)
	case "counter":
		value, err := strconv.ParseInt(rawValue, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h.storage.UpdateCounter(name, value)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")

	w.Header().Set("Content-Type", "text/plain")

	switch metricType {
	case "gauge":
		value, err := h.storage.GetGauge(name)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		fmt.Fprint(w, value)
	case "counter":
		value, err := h.storage.GetCounter(name)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		fmt.Fprint(w, value)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h *Handler) GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	html := "<html><body><h1>Metrics</h1><ul>"
	for name, value := range h.storage.GetAllGauges() {
		html += fmt.Sprintf("<li>%s (gauge) = %v</li>", name, value)
	}
	for name, value := range h.storage.GetAllCounters() {
		html += fmt.Sprintf("<li>%s (counter) = %v</li>", name, value)
	}
	html += "</ul></body></html>"

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, html)
}

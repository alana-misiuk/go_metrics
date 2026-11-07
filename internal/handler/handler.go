package handler

import (
	"bytes"
	"fmt"
	"go_metrics/internal/repository"
	"html/template"
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

	if name == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

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

var metricsPage = template.Must(template.New("metrics").Parse(`<html><body><h1>Metrics</h1><ul>
{{- range $name, $value := .Gauges}}<li>{{$name}} (gauge) = {{$value}}</li>{{end}}
{{- range $name, $value := .Counters}}<li>{{$name}} (counter) = {{$value}}</li>{{end -}}
</ul></body></html>`))

type metricsPageData struct {
	Gauges   map[string]float64
	Counters map[string]int64
}

func (h *Handler) GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	data := metricsPageData{
		Gauges:   h.storage.GetAllGauges(),
		Counters: h.storage.GetAllCounters(),
	}

	var buf bytes.Buffer
	if err := metricsPage.Execute(&buf, data); err != nil {
		http.Error(w, "failed to render metrics page", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w)
}

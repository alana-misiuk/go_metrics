package handler

import (
	"bytes"
	"errors"
	"fmt"
	models "go_metrics/internal/model"
	"go_metrics/internal/service"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type MetricsService interface {
	Update(metricType, name, rawValue string) error
	Get(metricType, name string) (models.Metrics, error)
	GetAll() (gauges map[string]float64, counters map[string]int64)
}

type Handler struct {
	svc MetricsService
}

func NewHandler(svc MetricsService) *Handler {
	return &Handler{svc: svc}
}

func statusFor(err error) int {
	switch {
	case errors.Is(err, service.ErrNameRequired), errors.Is(err, service.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, service.ErrUnknownType), errors.Is(err, service.ErrInvalidValue):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func (h *Handler) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	err := h.svc.Update(
		chi.URLParam(r, "type"),
		chi.URLParam(r, "name"),
		chi.URLParam(r, "value"),
	)
	if err != nil {
		w.WriteHeader(statusFor(err))
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetMetric(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	m, err := h.svc.Get(chi.URLParam(r, "type"), chi.URLParam(r, "name"))
	if err != nil {
		w.WriteHeader(statusFor(err))
		return
	}

	switch m.MType {
	case models.Gauge:
		fmt.Fprint(w, *m.Value)
	case models.Counter:
		fmt.Fprint(w, *m.Delta)
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
	gauges, counters := h.svc.GetAll()
	data := metricsPageData{Gauges: gauges, Counters: counters}

	var buf bytes.Buffer
	if err := metricsPage.Execute(&buf, data); err != nil {
		http.Error(w, "failed to render metrics page", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w)
}

package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	models "go_metrics/internal/model"
	"go_metrics/internal/service"
	"html/template"
	"net/http"

	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
)

type MetricsService interface {
	Update(metricType, name, rawValue string) error
	UpdateMetric(m models.Metrics) (models.Metrics, error)
	Get(metricType, name string) (models.Metrics, error)
	GetAll() (gauges map[string]float64, counters map[string]int64)
}

type Handler struct {
	svc      MetricsService
	log      *zap.Logger
	onUpdate func() error
	db       *sql.DB
}

func NewHandler(svc MetricsService, log *zap.Logger, onUpdate func() error, db *sql.DB) *Handler {
	return &Handler{svc: svc, log: log, onUpdate: onUpdate, db: db}
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := h.db.PingContext(r.Context()); err != nil {
		h.log.Error("database ping failed", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
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

func (h *Handler) writeError(w http.ResponseWriter, r *http.Request, err error) {
	status := statusFor(err)
	if status == http.StatusInternalServerError {
		h.log.Error("request error", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Error(err))
	}
	w.WriteHeader(status)
}

func (h *Handler) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	err := h.svc.Update(
		chi.URLParam(r, "type"),
		chi.URLParam(r, "name"),
		chi.URLParam(r, "value"),
	)
	if err != nil {
		h.writeError(w, r, err)
		return
	}

	if err := h.afterUpdate(); err != nil {
		h.writeError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetMetric(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	m, err := h.svc.Get(chi.URLParam(r, "type"), chi.URLParam(r, "name"))
	if err != nil {
		h.writeError(w, r, err)
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
		h.log.Error("render metrics page", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Error(err))
		http.Error(w, "failed to render metrics page", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w)
}

func (h *Handler) UpdateMetricJSON(w http.ResponseWriter, r *http.Request) {
	var req models.Metrics
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	m, err := h.svc.UpdateMetric(req)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	if err := h.afterUpdate(); err != nil {
		h.writeError(w, r, err)
		return
	}
	h.writeJSON(w, r, m)
}

func (h *Handler) afterUpdate() error {
	if h.onUpdate != nil {
		return h.onUpdate()
	}
	return nil
}

func (h *Handler) GetMetricJSON(w http.ResponseWriter, r *http.Request) {
	var req models.Metrics
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	m, err := h.svc.Get(req.MType, req.ID)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	h.writeJSON(w, r, m)
}

func (h *Handler) writeJSON(w http.ResponseWriter, r *http.Request, v any) {
	body, err := json.Marshal(v)
	if err != nil {
		h.log.Error("encode response", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

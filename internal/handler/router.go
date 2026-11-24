package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func NewRouter(h *Handler, log *zap.Logger) http.Handler {
	r := chi.NewRouter()
	r.Use(RequestLogger(log))
	r.Use(GzipMiddleware)

	r.Post("/update/{type}/{name}/{value}", h.UpdateMetric)
	r.Get("/value/{type}/{name}", h.GetMetric)
	r.Get("/", h.GetAllMetrics)
	r.Post("/update", h.UpdateMetricJSON)
	r.Post("/update/", h.UpdateMetricJSON)
	r.Post("/value", h.GetMetricJSON)
	r.Post("/value/", h.GetMetricJSON)
	return r
}

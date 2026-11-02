package handler

import (
	"go_metrics/internal/repository"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
    storage repository.Storage
}

func NewHandler(storage repository.Storage) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
	segments := strings.Split(r.URL.Path, "/")
	if len(segments) < 5 {
		w.WriteHeader(http.StatusNotFound)
        return
	}

	if segments[2] != "gauge" &&  segments[2] != "counter" {
		w.WriteHeader(http.StatusBadRequest)
        return
	}

	if segments[2] == "gauge" {
		value, err := strconv.ParseFloat(segments[4], 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
        	return
		}
		h.storage.UpdateGauge(segments[3], value)
	}

	if segments[2] == "counter" {
		value, err := strconv.ParseInt(segments[4], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h.storage.UpdateCounter(segments[3], value)
	}

	w.Header().Set("Content-Type", "text/plan")
	w.WriteHeader(http.StatusOK)
}

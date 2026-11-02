package main

import (
	"go_metrics/internal/handler"
	"go_metrics/internal/repository"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	storage := repository.NewMemStorage()
	h := handler.NewHandler(storage)

	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", h.UpdateMetric)
	r.Get("/value/{type}/{name}", h.GetMetric)
	r.Get("/", h.GetAllMetrics)

	http.ListenAndServe(":8080", r)
}

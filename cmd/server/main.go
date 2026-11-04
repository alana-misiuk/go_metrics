package main

import (
	"flag"
	"go_metrics/internal/handler"
	"go_metrics/internal/repository"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	addr := flag.String("a", "localhost:8080", "server address")
	flag.Parse()

	storage := repository.NewMemStorage()
	h := handler.NewHandler(storage)

	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", h.UpdateMetric)
	r.Get("/value/{type}/{name}", h.GetMetric)
	r.Get("/", h.GetAllMetrics)

	http.ListenAndServe(*addr, r)
}

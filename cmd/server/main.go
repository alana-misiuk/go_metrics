package main

import (
	"go_metrics/internal/handler"
	"go_metrics/internal/repository"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	storage := repository.NewMemStorage()
	handler:= handler.NewHandler(storage)
	mux.HandleFunc(`/update`, handler.UpdateMetric)
    http.ListenAndServe(`:8080`, mux)
}

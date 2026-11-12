package main

import (
	"errors"
	"flag"
	"fmt"
	"go_metrics/internal/config"
	"go_metrics/internal/handler"
	"go_metrics/internal/logger"
	"go_metrics/internal/repository"
	"go_metrics/internal/service"
	"log"
	"net/http"
	"os"
)

func main() {
	cfg, err := config.ParseServer(os.Args[1:])
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			os.Exit(0)
		}
		os.Exit(2)
	}

	if err := logger.Initialize(cfg.FlagLogLevel); err != nil {
		fmt.Fprintln(os.Stderr, "log level:", err)
		os.Exit(2)
	}

	storage := repository.NewMemStorage()
	svc := service.NewMetricsService(storage)
	h := handler.NewHandler(svc)
	r := handler.NewRouter(h)

	log.Fatal(http.ListenAndServe(cfg.Addr, r))
}

package main

import (
	"errors"
	"flag"
	"go_metrics/internal/config"
	"go_metrics/internal/handler"
	"go_metrics/internal/repository"
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

	storage := repository.NewMemStorage()
	h := handler.NewHandler(storage)

	r := handler.NewRouter(h)

	log.Fatal(http.ListenAndServe(cfg.Addr, r))
}

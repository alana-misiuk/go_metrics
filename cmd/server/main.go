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
	"os/signal"
	"syscall"
	"time"
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

	var fs *repository.FileStorage
	if cfg.FileStoragePath != "" {
		fs = repository.NewFileStorage(cfg.FileStoragePath, storage)

		if cfg.Restore {
			if err := fs.Restore(); err != nil {
				log.Printf("restore metrics: %v", err)
			}
		}
	}

	svc := service.NewMetricsService(storage)

	var syncSave func()
	if fs != nil && cfg.StoreInterval == 0 {
		syncSave = func() {
			if err := fs.Save(); err != nil {
				log.Printf("sync save: %v", err)
			}
		}
	}

	h := handler.NewHandler(svc, syncSave)
	r := handler.NewRouter(h)

	if fs != nil && cfg.StoreInterval > 0 {
		go func() {
			ticker := time.NewTicker(time.Duration(cfg.StoreInterval) * time.Second)
			defer ticker.Stop()
			for range ticker.C {
				if err := fs.Save(); err != nil {
					log.Printf("periodic save: %v", err)
				}
			}
		}()
	}

	if fs != nil {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigCh
			if err := fs.Save(); err != nil {
				log.Printf("save on shutdown: %v", err)
			}
			os.Exit(0)
		}()
	}

	log.Fatal(http.ListenAndServe(cfg.Addr, r))
}

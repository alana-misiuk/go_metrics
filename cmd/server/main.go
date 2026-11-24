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
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	cfg, err := config.ParseServer(os.Args[1:])
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			os.Exit(0)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	if err := logger.Initialize(cfg.FlagLogLevel); err != nil {
		fmt.Fprintln(os.Stderr, "log level:", err)
		os.Exit(2)
	}

	storage := repository.NewMemStorage()

	var fs *repository.FileStorage
	if cfg.FileStoragePath != "" {
		var err error
		fs, err = repository.NewFileStorage(cfg.FileStoragePath, storage, cfg.Restore)
		if err != nil {
			logger.Log.Error("restore metrics", zap.Error(err))
		}
	}

	svc := service.NewMetricsService(storage)

	var syncSave func() error
	if fs != nil && cfg.StoreInterval == 0 {
		syncSave = fs.Save
	}

	h := handler.NewHandler(svc, syncSave)
	r := handler.NewRouter(h)

	if fs != nil && cfg.StoreInterval > 0 {
		go func() {
			ticker := time.NewTicker(time.Duration(cfg.StoreInterval) * time.Second)
			defer ticker.Stop()
			for range ticker.C {
				if err := fs.Save(); err != nil {
					logger.Log.Error("periodic save", zap.Error(err))
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
				logger.Log.Error("save on shutdown", zap.Error(err))
			}
			os.Exit(0)
		}()
	}

	logger.Log.Fatal("server stopped", zap.Error(http.ListenAndServe(cfg.Addr, r)))
}

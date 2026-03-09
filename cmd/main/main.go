package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/config"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/logger"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/router"
)

func main() {
	log := logger.Init()

	cfg := config.Load()

	addr := ":" + cfg.Server.Port

	srv := &http.Server{
		Handler: router.New(cfg),
		Addr:    addr,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info("Server started", "host", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()

	<-stop
	log.Info("Shutting down server...")

	strContextTime := os.Getenv("CONTEXT_TIME")
	contextTime, err := strconv.Atoi(strContextTime)
	if err != nil {
		log.Warn("Context time was not found, use default context time - 5 seconds")
		contextTime = 5
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(contextTime)*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	log.Info("Server stopped gracefully")
}

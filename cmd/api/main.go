package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"test1/internal/app"
	"test1/internal/config"
	"test1/internal/database"
	"test1/internal/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db, err := database.NewPostgres(cfg.DBDSN)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer db.Close()

	application := app.New(cfg, db)
	application.Logger.Info("database connected")

	r := router.NewRouter(db, cfg, application.Logger)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	go func() {
		application.Logger.Info("starting server", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			application.Logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	application.Logger.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		application.Logger.Error("graceful shutdown failed", "error", err)
		if err := server.Close(); err != nil {
			application.Logger.Error("server close failed", "error", err)
		}
	}

	application.Logger.Info("server stopped")
}

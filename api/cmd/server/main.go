package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PetoAdam/homenavi-marketplace/api/internal/config"
	"github.com/PetoAdam/homenavi-marketplace/api/internal/db"
	"github.com/PetoAdam/homenavi-marketplace/api/internal/http/server"
)

func main() {
	cfg := config.Load()

	gormDB, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}
	defer func() {
		sqlDB, err := gormDB.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	}()

	if err := db.Migrate(context.Background(), gormDB); err != nil {
		log.Fatalf("db migrate failed: %v", err)
	}

	h := server.New(cfg, gormDB)

	srv := &http.Server{
		Addr:              cfg.BindAddress,
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("marketplace api listening on %s", cfg.BindAddress)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}
}

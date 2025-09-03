// Copyright 2025 Jacob Philip. All rights reserved.
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jacobtrvl/inventory-management/internal/api"
	"github.com/jacobtrvl/inventory-management/internal/inventory"
	"github.com/jacobtrvl/inventory-management/internal/store"
	"github.com/jacobtrvl/inventory-management/pkg/observability"
)

const (
	defaultAddr = ":8080"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	ctx := context.Background()
  
	db := store.NewMemDb()
	mc := observability.NewMetricsCollector()

	p := inventory.NewInventory(ctx, "products", db, mc)

	r := api.SetupRouter(ctx, p)
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = defaultAddr
	}
	srv := &http.Server{
		Addr:    addr,
		Handler: r.Handler(),
	}
	slog.Info("Starting server on ", "addr", srv.Addr)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server error", "error", err)
			close(stop)
		}
	}()
	<-stop
	slog.Info("Shutting down server...")
	mc.Shutdown()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown error", "error", err)
	}
	slog.Info("Server exiting")
}

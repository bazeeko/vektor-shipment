package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bazeeko/vektor-shipment/internal/app/grpc"
	shipmenthandler "github.com/bazeeko/vektor-shipment/internal/app/grpc/shipment"
	"github.com/bazeeko/vektor-shipment/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config.Load: %s", err)
	}

	shipmentHandler := shipmenthandler.New(cfg)

	server, err := grpc.NewServer(shipmentHandler, cfg.Server.GRPCPort)
	if err != nil {
		log.Fatalf("grpc.NewServer: %s", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err = server.Start(); err != nil {
			log.Fatalf("server.Start: %s", err)
		}
	}()

	<-quit

	slog.Info("Shutting down server...")

	server.GracefulStop()

	slog.Info("Successfully shut down server.")
}

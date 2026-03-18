package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bazeeko/vektor-shipment/internal/app/grpc"
	shipmenthandler "github.com/bazeeko/vektor-shipment/internal/app/grpc/shipment"
	"github.com/bazeeko/vektor-shipment/internal/config"
	"github.com/bazeeko/vektor-shipment/internal/repository/postgresql"
	shipmentrepository "github.com/bazeeko/vektor-shipment/internal/repository/postgresql/shipment"
	pricingservice "github.com/bazeeko/vektor-shipment/internal/services/pricing"
	"github.com/bazeeko/vektor-shipment/internal/services/reference"
	shipmentservice "github.com/bazeeko/vektor-shipment/internal/services/shipment"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config.Load: %s", err)
	}

	pool, err := postgresql.NewConnectWithMigration(ctx, postgresql.Params{
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		DBName:   cfg.Database.DBName,
	})
	if err != nil {
		log.Fatalf("postgresql.NewConnectWithMigration: %s", err)
	}

	referenceGenerator := reference.New()
	pricingService := pricingservice.New()

	shipmentRepository := shipmentrepository.New(pool)
	shipmentService := shipmentservice.New(shipmentRepository, referenceGenerator, pricingService)
	shipmentHandler := shipmenthandler.New(shipmentService)

	grpcServer, err := grpc.NewServer(shipmentHandler, cfg.Server.GRPCPort)
	if err != nil {
		log.Fatalf("grpc.NewServer: %s", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err = grpcServer.Start(); err != nil {
			log.Fatalf("grpcServer.Start: %s", err)
		}
	}()

	<-quit

	slog.Info("Shutting down grpc server...")

	grpcServer.GracefulStop()

	slog.Info("Successfully shut down grpc server.")
}

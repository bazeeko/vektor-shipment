package grpc

import (
	"fmt"
	"log/slog"
	"net"

	shipmentpb "github.com/bazeeko/vektor-shipment/pkg/api/shipment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	server *grpc.Server
	port   string
}

func NewServer(shipmentHandler shipmentpb.ShipmentServiceServer, port int) (*Server, error) {
	grpcServer := grpc.NewServer()

	reflection.Register(grpcServer)

	shipmentpb.RegisterShipmentServiceServer(grpcServer, shipmentHandler)

	return &Server{
		server: grpcServer,
		port:   fmt.Sprintf(":%d", port),
	}, nil
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.port)
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}

	slog.Info(fmt.Sprintf("Listening gRPC on %s", s.port))

	if err = s.server.Serve(listener); err != nil {
		return fmt.Errorf("s.server.Serve: %w", err)
	}

	return nil
}

func (s *Server) GracefulStop() {
	s.server.GracefulStop()
}

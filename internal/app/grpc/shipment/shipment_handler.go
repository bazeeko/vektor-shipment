package shipment

import (
	shipmentpb "github.com/bazeeko/vektor-shipment/pkg/api/shipment/v1"
)

type Service interface {
}

var _ shipmentpb.ShipmentServiceServer = (*Handler)(nil)

type Handler struct {
	shipmentpb.UnimplementedShipmentServiceServer
	shipmentService Service
}

func New(s Service) *Handler {
	return &Handler{
		shipmentService: s,
	}
}

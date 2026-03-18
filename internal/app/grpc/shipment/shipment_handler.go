package shipment

import (
	"context"
	"errors"

	"github.com/bazeeko/vektor-shipment/internal/models"
	"github.com/bazeeko/vektor-shipment/internal/models/errs"
	shipmentpb "github.com/bazeeko/vektor-shipment/pkg/api/shipment/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service interface {
	CreateShipment(ctx context.Context, request models.CreateShipmentRequest) (uuid.UUID, error)
	AddShipmentEvent(ctx context.Context, request models.AddShipmentEventRequest) error
	GetShipment(ctx context.Context, shipmentID uuid.UUID) (models.GetShipmentResponse, error)
	GetShipmentEvents(ctx context.Context, shipmentID uuid.UUID) ([]models.ShipmentEvent, error)
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

func (h *Handler) CreateShipment(ctx context.Context, in *shipmentpb.CreateShipmentRequest) (*shipmentpb.CreateShipmentResponse, error) {
	unitID, err := uuid.Parse(in.GetUnitId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	shipmentID, err := h.shipmentService.CreateShipment(ctx, models.CreateShipmentRequest{
		Origin:      in.Origin,
		Destination: in.Destination,
		DriverName:  in.DriverName,
		UnitID:      unitID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "h.shipmentService.CreateShipment: %v", err)
	}

	return &shipmentpb.CreateShipmentResponse{
		ShipmentId: shipmentID.String(),
	}, nil
}

func (h *Handler) GetShipment(ctx context.Context, in *shipmentpb.GetShipmentRequest) (*shipmentpb.GetShipmentResponse, error) {
	shipmentID, err := uuid.Parse(in.GetShipmentId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	shipment, err := h.shipmentService.GetShipment(ctx, shipmentID)
	switch {
	case errors.Is(err, errs.ErrShipmentNotFound):
		return nil, status.Error(codes.NotFound, err.Error())
	case err != nil:
		return nil, status.Errorf(codes.Internal, "h.shipmentService.GetShipment: %v", err)
	}

	return &shipmentpb.GetShipmentResponse{
		Shipment: &shipmentpb.Shipment{
			Id:              shipment.ID.String(),
			ReferenceNumber: shipment.ReferenceNumber,
			Origin:          shipment.Origin,
			Destination:     shipment.Destination,
			DriverName:      shipment.DriverName,
			UnitId:          shipment.UnitID.String(),
			ShipmentCost:    shipment.ShipmentCost,
			DriverRevenue:   shipment.DriverRevenue,
			CreatedAt:       timestamppb.New(shipment.CreatedAt),
			Status:          shipmentpb.ShipmentStatus(shipment.Status),
			UpdatedAt:       timestamppb.New(shipment.UpdatedAt),
		},
	}, nil
}

func (h *Handler) AddShipmentEvent(ctx context.Context, in *shipmentpb.AddShipmentEventRequest) (*emptypb.Empty, error) {
	shipmentID, err := uuid.Parse(in.GetShipmentId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	request := models.AddShipmentEventRequest{
		ShipmentID: shipmentID,
		Status:     in.GetStatus(),
		Details:    in.GetDetails(),
	}

	err = h.shipmentService.AddShipmentEvent(ctx, request)
	switch {
	case errors.Is(err, errs.ErrInvalidEventStatus):
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	case err != nil:
		return nil, status.Errorf(codes.Internal, "h.shipmentService.AddShipmentEvent: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (h *Handler) GetShipmentEvents(ctx context.Context, req *shipmentpb.GetShipmentEventsRequest) (*shipmentpb.GetShipmentEventsResponse, error) {
	shipmentID, err := uuid.Parse(req.GetShipmentId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	events, err := h.shipmentService.GetShipmentEvents(ctx, shipmentID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "h.shipmentService.GetShipmentEvents: %v", err)
	}

	protoEvents := make([]*shipmentpb.ShipmentEvent, 0, len(events))

	for i := range events {
		protoEvents = append(protoEvents, &shipmentpb.ShipmentEvent{
			Id:         events[i].ID.String(),
			ShipmentId: events[i].ShipmentID.String(),
			Status:     shipmentpb.ShipmentStatus(events[i].Status),
			Details:    events[i].Details,
			OccurredAt: timestamppb.New(events[i].OccurredAt),
		})
	}

	return &shipmentpb.GetShipmentEventsResponse{
		Events: protoEvents,
	}, nil
}

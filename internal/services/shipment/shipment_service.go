package shipment

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/bazeeko/vektor-shipment/internal/models"
	"github.com/bazeeko/vektor-shipment/internal/models/errs"
	shipmentrepository "github.com/bazeeko/vektor-shipment/internal/repository/postgresql/shipment"
	"github.com/google/uuid"
)

type Repository interface {
	InsertShipment(ctx context.Context, params shipmentrepository.InsertShipmentParams) (uuid.UUID, error)
	InsertEvent(ctx context.Context, params shipmentrepository.InsertEventParams) error
	SelectShipment(ctx context.Context, shipmentID uuid.UUID) (shipmentrepository.SelectShipmentOutput, error)
	SelectEvents(ctx context.Context, shipmentID uuid.UUID) ([]shipmentrepository.SelectEventOutput, error)
	SelectLastEvent(ctx context.Context, shipmentID uuid.UUID) (shipmentrepository.SelectEventOutput, error)
}

type ReferenceGenerator interface {
	GenerateReferenceNumber() (string, error)
}

type Service struct {
	shipmentRepository Repository
	referenceGenerator ReferenceGenerator
}

func New(shipmentRepository Repository, referenceGenerator ReferenceGenerator) *Service {
	return &Service{
		shipmentRepository: shipmentRepository,
		referenceGenerator: referenceGenerator,
	}
}

func (s *Service) CreateShipment(ctx context.Context, request models.CreateShipmentRequest) (uuid.UUID, error) {
	refNumber, err := s.referenceGenerator.GenerateReferenceNumber()
	if err != nil {
		return uuid.Nil, fmt.Errorf("s.referenceGenerator.GenerateReferenceNumber: %w", err)
	}

	insertShipmentParams := shipmentrepository.InsertShipmentParams{
		ReferenceNumber: refNumber,
		UnitID:          request.UnitID,
		Origin:          request.Origin,
		Destination:     request.Destination,
		DriverName:      request.DriverName,
		ShipmentCost:    request.ShipmentCost,
		DriverRevenue:   request.DriverRevenue,
		Status:          models.ShipmentStatusPending,
		EventDetails:    "Shipment created.",
	}

	shipmentID, err := s.shipmentRepository.InsertShipment(ctx, insertShipmentParams)
	if err != nil {
		return uuid.Nil, fmt.Errorf("s.shipmentRepository.InsertShipment: %w", err)
	}

	return shipmentID, nil
}

func (s *Service) AddShipmentEvent(ctx context.Context, request models.AddShipmentEventRequest) error {
	lastEvent, err := s.shipmentRepository.SelectLastEvent(ctx, request.ShipmentID)
	if err != nil {
		return fmt.Errorf("s.shipmentRepository.SelectLastEvent: %w", err)
	}

	isValidStatus := false

	switch lastEvent.Status {
	case models.ShipmentStatusPending:
		isValidStatus = request.Status == models.ShipmentStatusAwaitingDriver || request.Status == models.ShipmentStatusCancelled
	case models.ShipmentStatusAwaitingDriver:
		isValidStatus = request.Status == models.ShipmentStatusPickedUp || request.Status == models.ShipmentStatusCancelled
	case models.ShipmentStatusPickedUp:
		isValidStatus = request.Status == models.ShipmentStatusInTransit || request.Status == models.ShipmentStatusCancelled
	case models.ShipmentStatusInTransit:
		isValidStatus = request.Status == models.ShipmentStatusDelivered ||
			request.Status == models.ShipmentStatusDelayed ||
			request.Status == models.ShipmentStatusAtTransferPoint ||
			request.Status == models.ShipmentStatusCancelled
	case models.ShipmentStatusDelayed:
		isValidStatus = request.Status == models.ShipmentStatusInTransit || request.Status == models.ShipmentStatusCancelled
	case models.ShipmentStatusAtTransferPoint:
		isValidStatus = request.Status == models.ShipmentStatusAwaitingDriver || request.Status == models.ShipmentStatusCancelled
	case models.ShipmentStatusCancelled, models.ShipmentStatusDelivered:
		isValidStatus = false
	}

	if !isValidStatus {
		slog.Error(
			"invalid shipment event status",
			"shipmentID", request.ShipmentID,
			"currentStatus", lastEvent.Status,
			"newStatus", request.Status,
		)
		return errs.ErrInvalidEventStatus
	}

	insertEventParams := shipmentrepository.InsertEventParams{
		ShipmentID: request.ShipmentID,
		Status:     request.Status,
		Details:    request.Details,
	}

	if err = s.shipmentRepository.InsertEvent(ctx, insertEventParams); err != nil {
		return fmt.Errorf("s.shipmentRepository.InsertEvent: %w", err)
	}

	return nil
}

func (s *Service) GetShipment(ctx context.Context, shipmentID uuid.UUID) (models.GetShipmentResponse, error) {
	shipment, err := s.shipmentRepository.SelectShipment(ctx, shipmentID)
	if err != nil {
		return models.GetShipmentResponse{}, fmt.Errorf("s.shipmentRepository.SelectShipment: %w", err)
	}

	return models.GetShipmentResponse{
		ID:              shipment.ID,
		ReferenceNumber: shipment.ReferenceNumber,
		UnitID:          shipment.UnitID,
		Origin:          shipment.Origin,
		Destination:     shipment.Destination,
		DriverName:      shipment.DriverName,
		ShipmentCost:    shipment.ShipmentCost,
		DriverRevenue:   shipment.DriverRevenue,
		CreatedAt:       shipment.CreatedAt,
		Status:          shipment.Status,
		UpdatedAt:       shipment.UpdatedAt,
	}, nil
}

func (s *Service) GetShipmentEvents(ctx context.Context, shipmentID uuid.UUID) ([]models.ShipmentEvent, error) {
	events, err := s.shipmentRepository.SelectEvents(ctx, shipmentID)
	if err != nil {
		return nil, fmt.Errorf("s.shipmentRepository.SelectEvents: %w", err)
	}

	result := make([]models.ShipmentEvent, 0, len(events))

	for i := range events {
		result = append(result, models.ShipmentEvent{
			ID:         events[i].ID,
			ShipmentID: events[i].ShipmentID,
			Status:     events[i].Status,
			Details:    events[i].Details,
			OccurredAt: events[i].OccurredAt,
		})
	}

	return result, nil
}

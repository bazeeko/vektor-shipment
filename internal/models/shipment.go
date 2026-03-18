package models

import (
	"time"

	shipmentpb "github.com/bazeeko/vektor-shipment/pkg/api/shipment/v1"
	"github.com/google/uuid"
)

type ShipmentStatus int

const (
	ShipmentStatusUnknown ShipmentStatus = iota
	ShipmentStatusPending
	ShipmentStatusAwaitingDriver
	ShipmentStatusPickedUp
	ShipmentStatusInTransit
	ShipmentStatusDelayed
	ShipmentStatusAtTransferPoint
	ShipmentStatusDelivered
	ShipmentStatusCancelled
)

type CreateShipmentRequest struct {
	Origin        string
	Destination   string
	DriverName    string
	UnitID        uuid.UUID
	ShipmentCost  int64
	DriverRevenue int64
}

type GetShipmentResponse struct {
	ID              uuid.UUID
	ReferenceNumber string
	UnitID          uuid.UUID
	Origin          string
	Destination     string
	DriverName      string
	ShipmentCost    int64
	DriverRevenue   int64
	CreatedAt       time.Time
	Status          shipmentpb.ShipmentStatus
	UpdatedAt       time.Time
}

type AddShipmentEventRequest struct {
	ShipmentID uuid.UUID
	Status     shipmentpb.ShipmentStatus
	Details    string
}

type ShipmentEvent struct {
	ID         uuid.UUID
	ShipmentID uuid.UUID
	Status     shipmentpb.ShipmentStatus
	Details    string
	OccurredAt time.Time
}

package models

import (
	"time"

	"github.com/google/uuid"
)

type ShipmentStatus int

const (
	ShipmentStatusPending = iota
	ShipmentStatusAwaitingDriver
	ShipmentStatusPickedUp
	ShipmentStatusInTransit
	ShipmentStatusDelayed
	ShipmentStatusAtTransferPoint
	ShipmentStatusDelivered
	ShipmentStatusCancelled
)

type CreateShipmentRequest struct {
	ReferenceNumber string
	Origin          string
	Destination     string
	DriverName      string
	UnitID          uuid.UUID
	ShipmentCost    int64
	DriverRevenue   int64
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
}

type AddShipmentEventRequest struct {
	ShipmentID uuid.UUID
	Status     ShipmentStatus
	Details    string
}

type ShipmentEvent struct {
	ID         uuid.UUID
	ShipmentID uuid.UUID
	Status     ShipmentStatus
	Details    string
	OccurredAt time.Time
}

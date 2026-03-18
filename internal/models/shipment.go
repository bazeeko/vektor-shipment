package models

import (
	"time"

	"github.com/google/uuid"
)

type ShipmentStatus int

const (
	ShipmentStatusPending ShipmentStatus = iota
	ShipmentStatusAwaitingDriver
	ShipmentStatusPickedUp
	ShipmentStatusInTransit
	ShipmentStatusDelayed
	ShipmentStatusAtTransferPoint
	ShipmentStatusDelivered
	ShipmentStatusCancelled
)

func (s ShipmentStatus) String() string {
	switch s {
	case ShipmentStatusPending:
		return "PENDING"
	case ShipmentStatusAwaitingDriver:
		return "AWAITING_DRIVER"
	case ShipmentStatusPickedUp:
		return "PICKED_UP"
	case ShipmentStatusInTransit:
		return "IN_TRANSIT"
	case ShipmentStatusDelayed:
		return "DELAYED"
	case ShipmentStatusAtTransferPoint:
		return "AT_TRANSFER_POINT"
	case ShipmentStatusDelivered:
		return "DELIVERED"
	case ShipmentStatusCancelled:
		return "CANCELLED"
	default:
		return "UNKNOWN"
	}
}

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
	Status          ShipmentStatus
	UpdatedAt       time.Time
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

package postgresql

import (
	"time"

	shipmentpb "github.com/bazeeko/vektor-shipment/pkg/api/shipment/v1"
	"github.com/google/uuid"
)

type InsertShipmentParams struct {
	ReferenceNumber string
	UnitID          uuid.UUID
	Origin          string
	Destination     string
	DriverName      string
	ShipmentCost    int64
	DriverRevenue   int64
	Status          shipmentpb.ShipmentStatus
	EventDetails    string
}

type InsertEventParams struct {
	ShipmentID uuid.UUID
	Status     shipmentpb.ShipmentStatus
	Details    string
}

type SelectShipmentOutput struct {
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

type SelectEventOutput struct {
	ID         uuid.UUID
	ShipmentID uuid.UUID
	Status     shipmentpb.ShipmentStatus
	Details    string
	OccurredAt time.Time
}

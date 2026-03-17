package shipment

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	queryInsertShipment = `INSERT INTO public.shipments
    (reference_number, unit_id, origin, destination, driver_name, shipment_cost, driver_revenue)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING id;`

	queryInsertEvent = `INSERT INTO public.events
    (shipment_id, status, details)
	VALUES ($1, $2, $3)`

	querySelectShipment = `SELECT id, reference_number, unit_id, origin, destination, driver_name, shipment_cost, driver_revenue, created_at
	FROM public.shipments
	WHERE id = $1`

	querySelectEvents = `SELECT id, shipment_id, status, details, occurred_at
	FROM public.events
	WHERE shipment_id = $1`

	querySelectLastEvent = `SELECT id, shipment_id, status, details, occurred_at
	FROM public.events
	WHERE shipment_id = $1
	ORDER BY occurred_at DESC
	LIMIT 1`
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) InsertShipment(ctx context.Context, params InsertShipmentParams) error {
	var shipmentID uuid.UUID

	err := pgx.BeginTxFunc(ctx, r.pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		err := tx.QueryRow(
			ctx,
			queryInsertShipment,
			params.ReferenceNumber,
			params.UnitID,
			params.Origin,
			params.Destination,
			params.DriverName,
			params.ShipmentCost,
			params.DriverRevenue,
		).Scan(
			shipmentID,
		)
		if err != nil {
			return fmt.Errorf("tx.Exec queryInsertShipment: %w", err)
		}

		_, err = tx.Exec(
			ctx,
			queryInsertEvent,
			shipmentID,
			params.Status,
			params.EventDetails,
		)
		if err != nil {
			return fmt.Errorf("tx.Exec queryInsertEvent: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("pgx.BeginTxFunc: %w", err)
	}

	return nil
}

func (r *Repository) InsertEvent(ctx context.Context, params InsertEventParams) error {
	_, err := r.pool.Exec(
		ctx,
		queryInsertEvent,
		params.ShipmentID,
		params.Status,
		params.Details,
	)
	if err != nil {
		return fmt.Errorf("r.pool.Exec: %w", err)
	}

	return nil
}

func (r *Repository) SelectShipment(ctx context.Context, shipmentID uuid.UUID) (SelectShipmentOutput, error) {
	var output SelectShipmentOutput

	err := r.pool.QueryRow(
		ctx,
		querySelectShipment,
		shipmentID,
	).Scan(
		&output.ID,
		&output.ReferenceNumber,
		&output.UnitID,
		&output.Origin,
		&output.Destination,
		&output.DriverName,
		&output.ShipmentCost,
		&output.DriverRevenue,
		&output.CreatedAt,
	)
	if err != nil {
		return SelectShipmentOutput{}, fmt.Errorf("r.pool.QueryRow: %w", err)
	}

	return output, nil
}

func (r *Repository) SelectEvents(ctx context.Context, shipmentID uuid.UUID) ([]SelectEventOutput, error) {
	var events []SelectEventOutput

	rows, err := r.pool.Query(ctx, querySelectEvents, shipmentID)
	if err != nil {
		return nil, fmt.Errorf("r.pool.Query: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var output SelectEventOutput

		if err = rows.Scan(&output.ID, output.ShipmentID, &output.Status, &output.Details, output.OccurredAt); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		events = append(events, output)
	}

	return events, nil
}

func (r *Repository) SelectLastEvent(ctx context.Context, shipmentID uuid.UUID) (SelectEventOutput, error) {
	var output SelectEventOutput

	err := r.pool.QueryRow(
		ctx,
		querySelectLastEvent,
		shipmentID,
	).Scan(
		&output.ID,
		output.ShipmentID,
		&output.Status,
		&output.Details,
		output.OccurredAt,
	)
	if err != nil {
		return SelectEventOutput{}, fmt.Errorf("r.pool.QueryRow: %w", err)
	}

	return output, nil
}

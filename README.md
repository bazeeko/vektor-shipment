# Shipment Tracking Microservice

A simplified gRPC microservice for managing shipments and tracking their status changes throughout the transportation lifecycle.

## Overview

This service provides a reliable way to create shipments and record their status transitions (e.g., from `PENDING` to `PICKED_UP` to `DELIVERED`). It ensures that only valid status transitions are allowed according to business rules and maintains a complete history of all status changes as events.

## Architecture

The service follows a layered architecture to separate concerns and manage dependencies effectively:

- **Transport Layer (`internal/app/grpc`)**: Implements the gRPC server and handlers. It handles request validation and maps gRPC messages to internal models.
- **Service Layer (`internal/services`)**: Contains the core business logic. It orchestrates the flow of data between the handler and the repository, enforcing state machine rules for status transitions.
- **Repository Layer (`internal/repository`)**: Manages data persistence. The current implementation uses PostgreSQL for storing shipments and their event history.
- **Models (`internal/models`)**: Defines the core domain entities and shared types used across all layers.

## Design Decisions & Assumptions

### 1. Status State Machine
The shipment lifecycle is governed by a set of allowed status transitions:
- `PENDING` → `AWAITING_DRIVER` or `CANCELLED`
- `AWAITING_DRIVER` → `PICKED_UP` or `CANCELLED`
- `PICKED_UP` → `IN_TRANSIT` or `CANCELLED`
- `IN_TRANSIT` → `DELIVERED`, `DELAYED`, `AT_TRANSFER_POINT`, or `CANCELLED`
- `DELAYED` → `IN_TRANSIT` or `CANCELLED`
- `AT_TRANSFER_POINT` → `AWAITING_DRIVER` or `CANCELLED`
- `DELIVERED` and `CANCELLED` are terminal states.

### 2. Multi-hop Tracking
The system supports complex delivery routes through transfer points. A shipment can go from `IN_TRANSIT` to `AT_TRANSFER_POINT`, then back to `AWAITING_DRIVER` to be picked up by a new driver, forming a repeatable cycle.

### 3. Event-Driven History
Every status change is recorded as a separate event in the database. This provides a full audit trail and allows for easy retrieval of the shipment's history.

### 4. Database Migrations
PostgreSQL is used with `goose` for version-controlled database migrations, ensuring the schema is always consistent across environments.

## How to Run

### Prerequisites
- Docker and Docker Compose
- Go 1.26+

### Start the service
```bash
docker-compose up -d
```
This command starts the PostgreSQL database and the gRPC service. The service will automatically run migrations on startup.

### Stop the service
```bash
docker-compose down
```

## Testing

The service includes unit tests for the core business logic in the service layer.

### Run tests
```bash
go test ./... -v
```

### Mocking Strategy
The service layer uses `minimock` to mock the repository interface, allowing for isolated testing of business rules without requiring a live database.

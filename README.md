# Shipment Tracking Microservice

A simplified gRPC microservice for managing shipments and tracking their status changes throughout the transportation lifecycle.

## Overview

This service provides a reliable way to create shipments and record their status transitions. It ensures that only valid status transitions are allowed according to business rules and maintains a complete history of all status changes as events in a PostgreSQL database.

## Architecture

The service follows a clean, layered architecture:

- **Transport Layer (`internal/app/grpc`)**: Implements the gRPC server and handlers. Handles error mapping and request/response conversion.
- **Service Layer (`internal/services`)**: Contains the core business logic and orchestrates dependencies.
- **Repository Layer (`internal/repository`)**: Manages data persistence using PostgreSQL and structured SQL queries.
- **Models (`internal/models`)**: Defines domain entities and shared business types.

## Key Features & Design Decisions

### 1. Robust State Machine
Shipment statuses are managed via a controlled state machine to prevent illegal transitions (e.g., a `DELIVERED` shipment cannot become `CANCELLED`).
- **Initial Status**: `Pending`
- **Supported Transitions**: 
    - `Pending` → `AwaitingDriver`, `Cancelled`
    - `AwaitingDriver` → `PickedUp`, `Cancelled`
    - `PickedUp` → `InTransit`, `Cancelled`
    - `InTransit` → `Delivered`, `Delayed`, `AtTransferPoint`, `Cancelled`
    - `Delayed` → `InTransit`, `Cancelled`
    - `AtTransferPoint` → `AwaitingDriver`, `Cancelled`

### 2. Extensible Reference Generation
Shipment reference numbers are generated through a dedicated `ReferenceGenerator` interface, allowing for different generation strategies (e.g., random strings, sequential numbering, or UUIDs) without modifying the core service logic.

### 3. Event-Based Audit Trail
Every status change is persisted as a "Shipment Event," providing a full chronological history of the shipment's journey.

### 4. Dynamic Pricing Calculation
The `PricingService` handles the automated calculation of shipment costs and driver revenues during shipment creation. This decouples financial logic from the core tracking service.

### 5. Database Schema & Migrations
PostgreSQL is used for storage, with schema management handled by `goose` migrations located in the `migrations/` directory.

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

Run tests with:
```bash
go test -v ./...
```
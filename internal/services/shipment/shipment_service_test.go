package shipment

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bazeeko/vektor-shipment/internal/models"
	"github.com/bazeeko/vektor-shipment/internal/models/errs"
	shipmentrepository "github.com/bazeeko/vektor-shipment/internal/repository/postgresql/shipment"
	shipmentpb "github.com/bazeeko/vektor-shipment/pkg/api/shipment/v1"
	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	repositoryMock := NewRepositoryMock(t)
	referenceGeneratorMock := NewReferenceGeneratorMock(t)

	service := New(repositoryMock, referenceGeneratorMock)

	require.NotNil(t, service)
	require.NotNil(t, service.shipmentRepository)
	require.NotNil(t, service.referenceGenerator)
}

func TestService_CreateShipment(t *testing.T) {
	refNumber := "DH5Q3E0J5X2E"
	unitID := uuid.New()
	shipmentID := uuid.New()

	tests := []struct {
		name                   string
		request                models.CreateShipmentRequest
		expectedResponse       uuid.UUID
		mockSetup              func(repository *RepositoryMock, referenceGenerator *ReferenceGeneratorMock)
		expectedErrorSubstring string
	}{
		{
			name: "Success",
			request: models.CreateShipmentRequest{
				Origin:        "Origin",
				Destination:   "Destination",
				DriverName:    "Driver Name",
				UnitID:        unitID,
				ShipmentCost:  100000,
				DriverRevenue: 60000,
			},
			mockSetup: func(repository *RepositoryMock, referenceGenerator *ReferenceGeneratorMock) {
				referenceGenerator.GenerateReferenceNumberMock.
					Expect().
					Return(refNumber, nil)

				insertShipmentParams := shipmentrepository.InsertShipmentParams{
					ReferenceNumber: refNumber,
					UnitID:          unitID,
					Origin:          "Origin",
					Destination:     "Destination",
					DriverName:      "Driver Name",
					ShipmentCost:    100000,
					DriverRevenue:   60000,
					Status:          shipmentpb.ShipmentStatus_Pending,
					EventDetails:    "Shipment created.",
				}

				repository.InsertShipmentMock.
					Expect(minimock.AnyContext, insertShipmentParams).
					Return(shipmentID, nil)
			},
			expectedResponse:       shipmentID,
			expectedErrorSubstring: "",
		},
		{
			name: "Reference generator error",
			request: models.CreateShipmentRequest{
				Origin:        "Origin",
				Destination:   "Destination",
				DriverName:    "Driver Name",
				UnitID:        unitID,
				ShipmentCost:  100000,
				DriverRevenue: 60000,
			},
			mockSetup: func(_ *RepositoryMock, referenceGenerator *ReferenceGeneratorMock) {
				referenceGenerator.GenerateReferenceNumberMock.
					Expect().
					Return("", errors.New("reference generator error"))
			},
			expectedErrorSubstring: "reference generator error",
		},
		{
			name: "Repository error",
			request: models.CreateShipmentRequest{
				Origin:        "Origin",
				Destination:   "Destination",
				DriverName:    "Driver Name",
				UnitID:        unitID,
				ShipmentCost:  100000,
				DriverRevenue: 60000,
			},
			mockSetup: func(repository *RepositoryMock, referenceGenerator *ReferenceGeneratorMock) {
				referenceGenerator.GenerateReferenceNumberMock.
					Expect().
					Return(refNumber, nil)

				insertShipmentParams := shipmentrepository.InsertShipmentParams{
					ReferenceNumber: refNumber,
					UnitID:          unitID,
					Origin:          "Origin",
					Destination:     "Destination",
					DriverName:      "Driver Name",
					ShipmentCost:    100000,
					DriverRevenue:   60000,
					Status:          shipmentpb.ShipmentStatus_Pending,
					EventDetails:    "Shipment created.",
				}

				repository.InsertShipmentMock.
					Expect(minimock.AnyContext, insertShipmentParams).
					Return(uuid.Nil, errors.New("db error"))
			},
			expectedErrorSubstring: "db error",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := minimock.NewController(t)

			repositoryMock := NewRepositoryMock(c)
			referenceGeneratorMock := NewReferenceGeneratorMock(c)

			test.mockSetup(repositoryMock, referenceGeneratorMock)

			shipmentService := New(repositoryMock, referenceGeneratorMock)

			actualResponse, actualError := shipmentService.CreateShipment(context.Background(), test.request)
			if len(test.expectedErrorSubstring) > 0 {
				require.ErrorContains(t, actualError, test.expectedErrorSubstring)
			}

			require.Equal(t, test.expectedResponse, actualResponse)
		})
	}
}

func TestService_AddShipmentEvent(t *testing.T) {
	shipmentID := uuid.New()

	tests := []struct {
		name                   string
		request                models.AddShipmentEventRequest
		mockSetup              func(repository *RepositoryMock)
		expectedErrorSubstring string
	}{
		{
			name: "shipmentRepository.SelectLastEvent error",
			request: models.AddShipmentEventRequest{
				ShipmentID: shipmentID,
				Status:     shipmentpb.ShipmentStatus_AwaitingDriver,
			},
			mockSetup: func(repository *RepositoryMock) {
				repository.SelectLastEventMock.
					Expect(minimock.AnyContext, shipmentID).
					Return(shipmentrepository.SelectEventOutput{}, errors.New("mock error"))
			},
			expectedErrorSubstring: "mock error",
		},
		{
			name: "Pending->AwaitingDriver success",
			request: models.AddShipmentEventRequest{
				ShipmentID: shipmentID,
				Status:     shipmentpb.ShipmentStatus_AwaitingDriver,
			},
			mockSetup: func(repository *RepositoryMock) {
				lastEventMock := shipmentrepository.SelectEventOutput{
					ShipmentID: shipmentID,
					Status:     shipmentpb.ShipmentStatus_Pending,
				}

				repository.SelectLastEventMock.
					Expect(minimock.AnyContext, shipmentID).
					Return(lastEventMock, nil)

				insertEventParams := shipmentrepository.InsertEventParams{
					ShipmentID: shipmentID,
					Status:     shipmentpb.ShipmentStatus_AwaitingDriver,
				}

				repository.InsertEventMock.
					Expect(minimock.AnyContext, insertEventParams).
					Return(nil)
			},
			expectedErrorSubstring: "",
		},
		{
			name: "Pending->Cancelled success",
			request: models.AddShipmentEventRequest{
				ShipmentID: shipmentID,
				Status:     shipmentpb.ShipmentStatus_Cancelled,
			},
			mockSetup: func(repository *RepositoryMock) {
				lastEventMock := shipmentrepository.SelectEventOutput{
					ShipmentID: shipmentID,
					Status:     shipmentpb.ShipmentStatus_Pending,
				}

				repository.SelectLastEventMock.
					Expect(minimock.AnyContext, shipmentID).
					Return(lastEventMock, nil)

				insertEventParams := shipmentrepository.InsertEventParams{
					ShipmentID: shipmentID,
					Status:     shipmentpb.ShipmentStatus_Cancelled,
				}

				repository.InsertEventMock.
					Expect(minimock.AnyContext, insertEventParams).
					Return(nil)
			},
			expectedErrorSubstring: "",
		},
		{
			name: "PickedUp->InTransit success",
			request: models.AddShipmentEventRequest{
				ShipmentID: shipmentID,
				Status:     shipmentpb.ShipmentStatus_InTransit,
			},
			mockSetup: func(repository *RepositoryMock) {
				repository.SelectLastEventMock.
					Expect(minimock.AnyContext, shipmentID).
					Return(shipmentrepository.SelectEventOutput{Status: shipmentpb.ShipmentStatus_PickedUp}, nil)

				repository.InsertEventMock.
					Expect(minimock.AnyContext, shipmentrepository.InsertEventParams{
						ShipmentID: shipmentID,
						Status:     shipmentpb.ShipmentStatus_InTransit,
					}).
					Return(nil)
			},
		},
		{
			name: "Pending->PickedUp error",
			request: models.AddShipmentEventRequest{
				ShipmentID: shipmentID,
				Status:     shipmentpb.ShipmentStatus_PickedUp,
			},
			mockSetup: func(repository *RepositoryMock) {
				lastEventMock := shipmentrepository.SelectEventOutput{
					ShipmentID: shipmentID,
					Status:     shipmentpb.ShipmentStatus_Pending,
				}

				repository.SelectLastEventMock.
					Expect(minimock.AnyContext, shipmentID).
					Return(lastEventMock, nil)
			},
			expectedErrorSubstring: errs.ErrInvalidEventStatus.Error(),
		},
		{
			name: "InTransit->Delivered success",
			request: models.AddShipmentEventRequest{
				ShipmentID: shipmentID,
				Status:     shipmentpb.ShipmentStatus_Delivered,
			},
			mockSetup: func(repository *RepositoryMock) {
				repository.SelectLastEventMock.
					Expect(minimock.AnyContext, shipmentID).
					Return(shipmentrepository.SelectEventOutput{Status: shipmentpb.ShipmentStatus_InTransit}, nil)

				repository.InsertEventMock.
					Expect(minimock.AnyContext, shipmentrepository.InsertEventParams{
						ShipmentID: shipmentID,
						Status:     shipmentpb.ShipmentStatus_Delivered,
					}).
					Return(nil)
			},
		},
		{
			name: "InTransit->Delayed success",
			request: models.AddShipmentEventRequest{
				ShipmentID: shipmentID,
				Status:     shipmentpb.ShipmentStatus_Delayed,
			},
			mockSetup: func(repository *RepositoryMock) {
				repository.SelectLastEventMock.
					Expect(minimock.AnyContext, shipmentID).
					Return(shipmentrepository.SelectEventOutput{Status: shipmentpb.ShipmentStatus_InTransit}, nil)

				repository.InsertEventMock.
					Expect(minimock.AnyContext, shipmentrepository.InsertEventParams{
						ShipmentID: shipmentID,
						Status:     shipmentpb.ShipmentStatus_Delayed,
					}).
					Return(nil)
			},
		},
		{
			name: "InTransit->AtTransferPoint success",
			request: models.AddShipmentEventRequest{
				ShipmentID: shipmentID,
				Status:     shipmentpb.ShipmentStatus_AtTransferPoint,
			},
			mockSetup: func(repository *RepositoryMock) {
				repository.SelectLastEventMock.
					Expect(minimock.AnyContext, shipmentID).
					Return(shipmentrepository.SelectEventOutput{Status: shipmentpb.ShipmentStatus_InTransit}, nil)

				repository.InsertEventMock.
					Expect(minimock.AnyContext, shipmentrepository.InsertEventParams{
						ShipmentID: shipmentID,
						Status:     shipmentpb.ShipmentStatus_AtTransferPoint,
					}).
					Return(nil)
			},
		},
		{
			name: "Delayed->InTransit success",
			request: models.AddShipmentEventRequest{
				ShipmentID: shipmentID,
				Status:     shipmentpb.ShipmentStatus_InTransit,
			},
			mockSetup: func(repository *RepositoryMock) {
				repository.SelectLastEventMock.
					Expect(minimock.AnyContext, shipmentID).
					Return(shipmentrepository.SelectEventOutput{Status: shipmentpb.ShipmentStatus_Delayed}, nil)

				repository.InsertEventMock.
					Expect(minimock.AnyContext, shipmentrepository.InsertEventParams{
						ShipmentID: shipmentID,
						Status:     shipmentpb.ShipmentStatus_InTransit,
					}).
					Return(nil)
			},
		},
		{
			name: "AtTransferPoint->AwaitingDriver success",
			request: models.AddShipmentEventRequest{
				ShipmentID: shipmentID,
				Status:     shipmentpb.ShipmentStatus_AwaitingDriver,
			},
			mockSetup: func(repository *RepositoryMock) {
				repository.SelectLastEventMock.
					Expect(minimock.AnyContext, shipmentID).
					Return(shipmentrepository.SelectEventOutput{Status: shipmentpb.ShipmentStatus_AtTransferPoint}, nil)

				repository.InsertEventMock.
					Expect(minimock.AnyContext, shipmentrepository.InsertEventParams{
						ShipmentID: shipmentID,
						Status:     shipmentpb.ShipmentStatus_AwaitingDriver,
					}).
					Return(nil)
			},
		},
		{
			name: "Delivered->Cancelled error",
			request: models.AddShipmentEventRequest{
				ShipmentID: shipmentID,
				Status:     shipmentpb.ShipmentStatus_Cancelled,
			},
			mockSetup: func(repository *RepositoryMock) {
				repository.SelectLastEventMock.
					Expect(minimock.AnyContext, shipmentID).
					Return(shipmentrepository.SelectEventOutput{Status: shipmentpb.ShipmentStatus_Delivered}, nil)
			},
			expectedErrorSubstring: errs.ErrInvalidEventStatus.Error(),
		},
		{
			name: "InsertEvent error",
			request: models.AddShipmentEventRequest{
				ShipmentID: shipmentID,
				Status:     shipmentpb.ShipmentStatus_Cancelled,
			},
			mockSetup: func(repository *RepositoryMock) {
				repository.SelectLastEventMock.
					Expect(minimock.AnyContext, shipmentID).
					Return(shipmentrepository.SelectEventOutput{Status: shipmentpb.ShipmentStatus_Pending}, nil)

				repository.InsertEventMock.
					Expect(minimock.AnyContext, shipmentrepository.InsertEventParams{
						ShipmentID: shipmentID,
						Status:     shipmentpb.ShipmentStatus_Cancelled,
					}).
					Return(errors.New("insert error"))
			},
			expectedErrorSubstring: "insert error",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := minimock.NewController(t)

			repositoryMock := NewRepositoryMock(c)
			referenceGeneratorMock := NewReferenceGeneratorMock(c)

			test.mockSetup(repositoryMock)

			shipmentService := New(repositoryMock, referenceGeneratorMock)

			actualError := shipmentService.AddShipmentEvent(context.Background(), test.request)
			if len(test.expectedErrorSubstring) > 0 {
				require.ErrorContains(t, actualError, test.expectedErrorSubstring)
			}
		})
	}
}

func TestService_GetShipment(t *testing.T) {
	mockUUID := uuid.New()

	tests := []struct {
		name                   string
		request                uuid.UUID
		expectedResponse       models.GetShipmentResponse
		mockSetup              func(repository *RepositoryMock)
		expectedErrorSubstring string
	}{
		{
			name:    "Shipment not found",
			request: mockUUID,
			mockSetup: func(repository *RepositoryMock) {
				repository.SelectShipmentMock.
					Expect(minimock.AnyContext, mockUUID).
					Return(shipmentrepository.SelectShipmentOutput{}, errs.ErrShipmentNotFound)
			},
			expectedResponse:       models.GetShipmentResponse{},
			expectedErrorSubstring: errs.ErrShipmentNotFound.Error(),
		},
		{
			name:    "Success",
			request: mockUUID,
			mockSetup: func(repository *RepositoryMock) {
				selectShipmentOutputMock := shipmentrepository.SelectShipmentOutput{
					ID:              mockUUID,
					ReferenceNumber: "refNumber",
					UnitID:          mockUUID,
					Origin:          "origin",
					Destination:     "destination",
					DriverName:      "driverName",
					ShipmentCost:    10000,
					DriverRevenue:   9999,
					CreatedAt:       time.Time{},
					Status:          shipmentpb.ShipmentStatus_Pending,
					UpdatedAt:       time.Time{},
				}

				repository.SelectShipmentMock.
					Expect(minimock.AnyContext, mockUUID).
					Return(selectShipmentOutputMock, nil)
			},
			expectedResponse: models.GetShipmentResponse{
				ID:              mockUUID,
				ReferenceNumber: "refNumber",
				UnitID:          mockUUID,
				Origin:          "origin",
				Destination:     "destination",
				DriverName:      "driverName",
				ShipmentCost:    10000,
				DriverRevenue:   9999,
				CreatedAt:       time.Time{},
				Status:          shipmentpb.ShipmentStatus_Pending,
				UpdatedAt:       time.Time{},
			},
			expectedErrorSubstring: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := minimock.NewController(t)

			repositoryMock := NewRepositoryMock(c)
			referenceGeneratorMock := NewReferenceGeneratorMock(c)

			test.mockSetup(repositoryMock)

			shipmentService := New(repositoryMock, referenceGeneratorMock)

			actualResponse, actualError := shipmentService.GetShipment(context.Background(), test.request)
			if len(test.expectedErrorSubstring) > 0 {
				require.ErrorContains(t, actualError, test.expectedErrorSubstring)
			}

			require.Equal(t, test.expectedResponse, actualResponse)
		})
	}
}

func TestService_GetShipmentEvents(t *testing.T) {
	mockUUID := uuid.New()

	tests := []struct {
		name                   string
		request                uuid.UUID
		expectedResponse       []models.ShipmentEvent
		mockSetup              func(repository *RepositoryMock)
		expectedErrorSubstring string
	}{
		{
			name:    "Shipment not found",
			request: mockUUID,
			mockSetup: func(repository *RepositoryMock) {
				repository.SelectEventsMock.
					Expect(minimock.AnyContext, mockUUID).
					Return([]shipmentrepository.SelectEventOutput{}, errs.ErrShipmentNotFound)
			},
			expectedResponse:       nil,
			expectedErrorSubstring: errs.ErrShipmentNotFound.Error(),
		},
		{
			name:    "Success",
			request: mockUUID,
			mockSetup: func(repository *RepositoryMock) {
				selectEventsOutputMock := []shipmentrepository.SelectEventOutput{
					{
						ID:         mockUUID,
						ShipmentID: mockUUID,
						Status:     shipmentpb.ShipmentStatus_Pending,
						Details:    "",
						OccurredAt: time.Time{},
					},
					{
						ID:         mockUUID,
						ShipmentID: mockUUID,
						Status:     shipmentpb.ShipmentStatus_AwaitingDriver,
						Details:    "",
						OccurredAt: time.Time{},
					},
					{
						ID:         mockUUID,
						ShipmentID: mockUUID,
						Status:     shipmentpb.ShipmentStatus_PickedUp,
						Details:    "",
						OccurredAt: time.Time{},
					},
				}

				repository.SelectEventsMock.
					Expect(minimock.AnyContext, mockUUID).
					Return(selectEventsOutputMock, nil)
			},
			expectedResponse: []models.ShipmentEvent{
				{
					ID:         mockUUID,
					ShipmentID: mockUUID,
					Status:     shipmentpb.ShipmentStatus_Pending,
					Details:    "",
					OccurredAt: time.Time{},
				},
				{
					ID:         mockUUID,
					ShipmentID: mockUUID,
					Status:     shipmentpb.ShipmentStatus_AwaitingDriver,
					Details:    "",
					OccurredAt: time.Time{},
				},
				{
					ID:         mockUUID,
					ShipmentID: mockUUID,
					Status:     shipmentpb.ShipmentStatus_PickedUp,
					Details:    "",
					OccurredAt: time.Time{},
				},
			},
			expectedErrorSubstring: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := minimock.NewController(t)

			repositoryMock := NewRepositoryMock(c)
			referenceGeneratorMock := NewReferenceGeneratorMock(c)

			test.mockSetup(repositoryMock)

			shipmentService := New(repositoryMock, referenceGeneratorMock)

			actualResponse, actualError := shipmentService.GetShipmentEvents(context.Background(), test.request)
			if len(test.expectedErrorSubstring) > 0 {
				require.ErrorContains(t, actualError, test.expectedErrorSubstring)
			}

			require.Equal(t, test.expectedResponse, actualResponse)
		})
	}
}

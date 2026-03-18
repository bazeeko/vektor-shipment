package pricing

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type Cost struct {
	ShipmentCost  int64
	DriverRevenue int64
}

type Service struct {
	maxCost *big.Int
}

func New() *Service {
	return &Service{
		maxCost: big.NewInt(100000),
	}
}

func (s *Service) CalculateShipmentCost() (Cost, error) {
	shipmentCost, err := rand.Int(rand.Reader, s.maxCost)
	if err != nil {
		return Cost{}, fmt.Errorf("rand.Int: %v", err)
	}

	shipmentCost.Sub(shipmentCost, new(big.Int).Mod(shipmentCost, big.NewInt(100)))

	driverRevenue := big.NewInt(shipmentCost.Int64())
	driverRevenue.Mul(driverRevenue, big.NewInt(80)).Div(driverRevenue, big.NewInt(100))

	return Cost{
		ShipmentCost:  shipmentCost.Int64(),
		DriverRevenue: driverRevenue.Int64(),
	}, nil
}

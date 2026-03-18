package errs

import "errors"

var (
	ErrInvalidEventStatus = errors.New("invalid shipment event status")
	ErrShipmentNotFound   = errors.New("shipment not found")
)

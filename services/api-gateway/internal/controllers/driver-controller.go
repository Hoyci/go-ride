package controllers

import (
	pd "go-ride/shared/proto/driver"

	"github.com/go-playground/validator/v10"
)

type DriverController struct {
	validator     *validator.Validate
	driverService pd.DriverServiceClient
}

func NewDriverController(v *validator.Validate, ts pd.DriverServiceClient) *DriverController {
	return &DriverController{
		validator:     v,
		driverService: ts,
	}
}

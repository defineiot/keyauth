package domain

import (
	"sync"

	"openauth/api/exception"
	"openauth/api/logger"
	"openauth/store/domain"
)

var (
	controller *Controller
	once       sync.Once
)

// GetController use to new an controller
func GetController() (*Controller, error) {
	if controller == nil {
		return nil, exception.NewInternalServerError("domain controller not initial")
	}
	return controller, nil
}

// InitController use to init controller
func InitController(ds domain.Store, log logger.OpenAuthLogger) {
	once.Do(func() {
		controller = &Controller{ds: ds, log: log}
		controller.log.Debug("initial domain controller successful")
	})
	controller.log.Info("domain controller already initialed")
}

// Controller is domain pkg
type Controller struct {
	ds  domain.Store
	log logger.OpenAuthLogger
}

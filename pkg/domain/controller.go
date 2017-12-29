package domain

import (
	"openauth/api/logger"
	"openauth/store/domain"
)

// NewController use to init controller
func NewController(ds domain.Store, log logger.OpenAuthLogger) *Controller {
	return &Controller{ds: ds, log: log}
}

// Controller is domain pkg
type Controller struct {
	ds  domain.Store
	log logger.OpenAuthLogger
}

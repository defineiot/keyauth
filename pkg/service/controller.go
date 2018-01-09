package service

import (
	"openauth/api/logger"
	"openauth/store/service"
)

// NewController use to init controller
func NewController(ss service.Store, log logger.OpenAuthLogger) *Controller {
	return &Controller{ss: ss, log: log}
}

// Controller is domain pkg
type Controller struct {
	ss  service.Store
	log logger.OpenAuthLogger
}

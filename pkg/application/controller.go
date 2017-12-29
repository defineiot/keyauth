package application

import (
	"openauth/api/logger"
	"openauth/store/application"
	"openauth/store/user"
)

// NewController use to initial controller
func NewController(logger logger.OpenAuthLogger, as application.Store, us user.Store) *Controller {
	controller := &Controller{log: logger, as: as, us: us}
	return controller
}

// Controller is domain pkg
type Controller struct {
	log logger.OpenAuthLogger
	as  application.Store
	us  user.Store
}

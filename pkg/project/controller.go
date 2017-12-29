package project

import (
	"openauth/api/logger"
	"openauth/store/domain"
	"openauth/store/project"
	"openauth/store/user"
)

// NewController use to initial controller
func NewController(logger logger.OpenAuthLogger, ds domain.Store, ps project.Store, us user.Store) *Controller {
	return &Controller{log: logger, ds: ds, ps: ps, us: us}
}

// Controller is domain pkg
type Controller struct {
	log logger.OpenAuthLogger
	ps  project.Store
	ds  domain.Store
	us  user.Store
}

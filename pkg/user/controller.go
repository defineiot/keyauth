package user

import (
	"openauth/api/logger"
	"openauth/store/domain"
	"openauth/store/project"
	"openauth/store/user"
)

// NewController use to initial user controller
func NewController(log logger.OpenAuthLogger, us user.Store, ds domain.Store, ps project.Store) *Controller {
	return &Controller{ds: ds, ps: ps, us: us, log: log}

}

// Controller is domain pkg
type Controller struct {
	log logger.OpenAuthLogger
	ps  project.Store
	ds  domain.Store
	us  user.Store
}

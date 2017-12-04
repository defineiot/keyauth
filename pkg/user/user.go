package user

import (
	"openauth/api/exception"
	"openauth/api/logger"
	"openauth/storage/domain"
	"openauth/storage/project"
	"openauth/storage/user"
)

// NewController use to new an controller
func NewController(logger logger.OpenAuthLogger, us user.Storage, ds domain.Storage, ps project.Storage) *Controller {
	return &Controller{logger: logger, ds: ds, ps: ps, us: us}
}

// Controller is domain pkg
type Controller struct {
	logger logger.OpenAuthLogger
	ps     project.Storage
	ds     domain.Storage
	us     user.Storage
}

// CreateUser create user
func (c *Controller) CreateUser(domainID, userName, password, description string) (*user.User, error) {
	ok, err := c.ds.CheckDomainIsExistByID(domainID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, exception.NewBadRequest("domain %s not exist", domainID)
	}

	u, err := c.us.CreateUser(domainID, userName, password, true, 4096, 4096)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// GetUser get on user
func (c *Controller) GetUser(userID string) (*user.User, error) {
	c.logger.Debugf("user id %s", userID)
	u, err := c.us.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// ListUsers list user
func (c *Controller) ListUsers(domainID string) ([]*user.User, error) {

	return nil, nil
}

// DeleteUser delete user
func (c *Controller) DeleteUser() error {
	return nil
}

// ListUserProjects list all projects
func (c *Controller) ListUserProjects() ([]*project.Project, error) {
	return nil, nil
}

// AddProjects add projects
func (c *Controller) AddProjects() error {
	return nil
}

// RemoveProjects remove projects
func (c *Controller) RemoveProjects() error {
	return nil
}

package project

import (
	"openauth/api/exception"
	"openauth/api/logger"
	"openauth/storage/domain"
	"openauth/storage/project"
	"openauth/storage/user"
)

// NewController use to new an controller
func NewController(logger logger.OpenAuthLogger, ds domain.Storage, ps project.Storage) *Controller {
	return &Controller{logger: logger, ds: ds, ps: ps}
}

// Controller is domain pkg
type Controller struct {
	logger logger.OpenAuthLogger
	ps     project.Storage
	ds     domain.Storage
}

// CreateProject use to create domain
func (c *Controller) CreateProject(domainID, name, description string, cred user.Credential) (*project.Project, error) {
	ok, err := c.ds.CheckDomainIsExistByID(domainID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, exception.NewBadRequest("domain %s not exist", domainID)
	}

	proj, err := c.ps.CreateProject(domainID, name, description, true)
	if err != nil {
		return nil, err
	}

	return proj, nil
}

// ListProject list an domain's project
func (c *Controller) ListProject(domainID string) ([]*project.Project, error) {
	ok, err := c.ds.CheckDomainIsExistByID(domainID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, exception.NewBadRequest("domain %s not exist", domainID)
	}

	projects, err := c.ps.ListDomainProjects(domainID)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

// GetProject use to get one project
func (c *Controller) GetProject(id string, cred user.Credential) (*project.Project, error) {
	proj, err := c.ps.GetProject(id)
	if err != nil {
		return nil, err
	}

	// TODO: check the project is for this user

	return proj, nil
}

// UpdateProject use to update one project
func (c *Controller) UpdateProject(cred user.Credential) (*project.Project, error) {
	return nil, nil
}

// DestroyProject use to delete one project
func (c *Controller) DestroyProject(id string, cred user.Credential) error {
	// TODO: check the projcet is for this user

	if err := c.ps.DeleteProject(id); err != nil {
		return err
	}

	return nil
}

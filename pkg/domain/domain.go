package domain

import (
	"openauth/api/logger"
	"openauth/storage/domain"
	"openauth/storage/user"
)

// NewController use to new an controller
func NewController(logger logger.OpenAuthLogger, ds domain.Storage) *Controller {
	return &Controller{logger: logger, ds: ds}
}

// Controller is domain pkg
type Controller struct {
	logger logger.OpenAuthLogger
	ds     domain.Storage
}

// CreateDomain use to create domain
func (c *Controller) CreateDomain(name, description, displayName string, cert user.Credential) (*domain.Domain, error) {
	dom, err := c.ds.CreateDomain(name, description, displayName, true)
	if err != nil {
		return nil, err
	}

	return dom, nil
}

// ListDomain use to list all domains
func (c *Controller) ListDomain() ([]*domain.Domain, error) {
	doms, err := c.ds.ListDomain()
	if err != nil {
		return nil, err
	}

	return doms, nil
}

// GetDomain use to get an domain
func (c *Controller) GetDomain(domainID string) (*domain.Domain, error) {
	dom, err := c.ds.GetDomain(domainID)
	if err != nil {
		return nil, err
	}

	return dom, nil

}

// UpdateDomain use to update an domain
func (c *Controller) UpdateDomain() {

}

// DestoryDomain use to delete an domain
func (c *Controller) DestoryDomain(domainID string) error {
	if err := c.ds.DeleteDomain(domainID); err != nil {
		return err
	}

	return nil
}

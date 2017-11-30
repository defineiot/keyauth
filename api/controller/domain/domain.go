package domain

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/sirupsen/logrus"

	"openauth/api/exception"
	"openauth/pkg/domain"
	"openauth/pkg/user"
)

var (
	controller *Controller
	once       sync.Once
)

// GetController use to use an domain controller
func GetController() (*Controller, error) {
	if controller == nil {
		return nil, errors.New("domain controller isn't initial")
	}

	return controller, nil
}

// InitController use to initial an domain controller instance
func InitController(db *sql.DB, logger *logrus.Logger, domain domain.Manager) error {
	once.Do(func() {
		controller = &Controller{db: db, logger: logger, dm: domain}
	})

	return nil
}

// Controller is domain controller
type Controller struct {
	db     *sql.DB
	logger *logrus.Logger
	dm     domain.Manager
}

// CreateDomain use to create domain
func (c *Controller) CreateDomain(name, description, displayName string, cert user.Credential) (*domain.Domain, error) {
	dom, err := c.dm.CreateDomain(name, description, displayName, false)
	if err != nil {
		return nil, fmt.Errorf("create domain error, %s", err)
	}

	return dom, nil
}

// ListDomain use to list all domains
func (c *Controller) ListDomain() ([]*domain.Domain, error) {
	doms, err := c.dm.ListDomain()
	if err != nil {
		return nil, fmt.Errorf("list domain error, %s", err)
	}

	return doms, nil
}

// GetDomain use to get an domain
func (c *Controller) GetDomain(domainID string) (*domain.Domain, error) {
	dom, err := c.dm.GetDomain(domainID)
	if err != nil {
		return nil, err
	}
	if dom == nil {
		return nil, exception.NewAPIException(fmt.Sprintf("domain %s not find", domainID), http.StatusNotFound)
	}

	return dom, nil

}

// UpdateDomain use to update an domain
func (c *Controller) UpdateDomain() {

}

// DestoryDomain use to delete an domain
func (c *Controller) DestoryDomain(domainID string) error {
	if err := c.dm.DeleteDomain(domainID); err != nil {
		return fmt.Errorf("delete domain error, %s", err)
	}

	return nil
}

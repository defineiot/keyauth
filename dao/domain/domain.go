package domain

import (
	"github.com/defineiot/keyauth/dao/models"
)

// Store is an domain service
type Store interface {
	Reader
	Writer
	Close() error
}

// Reader for read data from store
type Reader interface {
	GetDomainByID(domainID string) (*models.Domain, error)
	GetDomainByName(name string) (*models.Domain, error)
	CheckDomainIsExistByID(domainID string) (bool, error)
	CheckDomainIsExistByName(domainName string) (bool, error)
	ListDomain(pageNumber, pageSize int64) (domains []*models.Domain, totalPage int64, err error)
	ListUserThirdDomains(userID string) ([]*models.Domain, error)
}

// Writer for write data to store
type Writer interface {
	CreateDomain(d *models.Domain) error
	UpdateDomain(id, name, description string) (*models.Domain, error)
	DeleteDomainByID(id string) error
	DeleteDomainByName(name string) error
}

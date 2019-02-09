package service

import (
	"github.com/defineiot/keyauth/dao/models"
)

// Store service store interface
type Store interface {
	Reader
	Writer
	Close() error
}

// Reader read service information from store
type Reader interface {
	ListServices() ([]*models.Service, error)
	GetServiceByID(id string) (*models.Service, error)
	GetServiceByName(name string) (*models.Service, error)
	GetServiceByClientID(clientID string) (*models.Service, error)

	ListAllFeatures() ([]*models.Feature, error)
	ListServiceFeatures(serviceID string) ([]*models.Feature, error)
	ListRoleFeatures(roleID string) ([]*models.Feature, error)
}

// Writer write service information to store
type Writer interface {
	CreateService(service *models.Service) error
	DeleteService(id string) error

	RegistryServiceFeatures(serviceID, version string, features ...*models.Feature) error
	AssociateFeaturesToRole(roleID string, features ...*models.Feature) error
	UnlinkFeatureFromRole(roleID string, features ...*models.Feature) error
}

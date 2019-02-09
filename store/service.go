package store

import (
	"github.com/defineiot/keyauth/dao/models"
	"github.com/defineiot/keyauth/internal/exception"
)

// CreateService todo
func (s *Store) CreateService(svr *models.Service) error {
	return s.dao.Service.CreateService(svr)
}

// ListServices todo
func (s *Store) ListServices() ([]*models.Service, error) {
	return s.dao.Service.ListServices()
}

// GetService todo
func (s *Store) GetService(id string) (*models.Service, error) {
	return s.dao.Service.GetServiceByID(id)
}

// GetServiceByName todo
func (s *Store) GetServiceByName(name string) (*models.Service, error) {
	return s.dao.Service.GetServiceByName(name)
}

// DeleteService todo
func (s *Store) DeleteService(id string) error {
	return s.dao.Service.DeleteService(id)
}

// RegistryServiceFeatures todo
func (s *Store) RegistryServiceFeatures(id, version string, features ...*models.Feature) error {
	if version == "" {
		return exception.NewBadRequest("service version is \"\"")
	}

	for i := range features {
		if err := features[i].Validate(); err != nil {
			return err
		}
	}

	return s.dao.Service.RegistryServiceFeatures(id, version, features...)
}

// ListServiceFeatures todo
func (s *Store) ListServiceFeatures(name string) ([]*models.Feature, error) {
	return s.dao.Service.ListServiceFeatures(name)
}

// CheckServiceHasFeature todo
func (s *Store) CheckServiceHasFeature(sn, fn string) (bool, error) {
	svr, err := s.dao.Service.GetServiceByName(sn)
	if err != nil {
		return false, err
	}

	var exist bool
	for i := range svr.Features {
		if svr.Features[i].Name == fn {
			exist = true
			break
		}
	}

	return exist, nil
}

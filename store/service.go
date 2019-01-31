package store

import (
	"github.com/defineiot/keyauth/dao/service"
)

// CreateService todo
func (s *Store) CreateService(svr *service.Service) error {
	return s.dao.Service.CreateService(svr)
}

// ListServices todo
func (s *Store) ListServices() ([]*service.Service, error) {
	return s.dao.Service.ListServices()
}

// GetService todo
func (s *Store) GetService(id string) (*service.Service, error) {
	return s.dao.Service.GetServiceByID(id)
}

// DeleteService todo
func (s *Store) DeleteService(id string) error {
	return s.dao.Service.DeleteService(id)
}

// RegistryServiceFeatures todo
func (s *Store) RegistryServiceFeatures(id, version string, features ...*service.Feature) error {
	for i := range features {
		if err := features[i].Validate(); err != nil {
			return err
		}
	}

	return s.dao.Service.RegistryServiceFeatures(id, version, features...)
}

// ListServiceFeatures todo
func (s *Store) ListServiceFeatures(name string) ([]*service.Feature, error) {
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

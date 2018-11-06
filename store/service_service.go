package store

import (
	"github.com/defineiot/keyauth/dao/client"
	"github.com/defineiot/keyauth/dao/service"
	"github.com/defineiot/keyauth/internal/exception"
)

// CreateService todo
func (s *Store) CreateService(name, description string) (*service.Service, error) {
	cli, err := s.client.CreateClient(client.Confidential, "")
	if err != nil {
		return nil, err
	}

	svr, err := s.service.CreateService(name, description, cli.ID)
	if err != nil {
		return nil, err
	}

	svr.Client = cli

	return svr, nil
}

// CheckService todo
func (s *Store) CheckService(name string) (bool, error) {
	return s.service.CheckServiceIsExist(name)
}

// ListServices todo
func (s *Store) ListServices() ([]*service.Service, error) {
	svrs, err := s.service.ListServices()
	if err != nil {
		return nil, err
	}

	for _, svr := range svrs {
		if svr.ClientID != "" {
			cli, err := s.client.GetClient(svr.ClientID)
			if err != nil {
				return nil, err
			}
			svr.Client = cli
		}
	}

	return svrs, nil
}

// GetService todo
func (s *Store) GetService(name string) (*service.Service, error) {
	svr, err := s.service.GetService(name)
	if err != nil {
		return nil, err
	}

	if svr.ClientID != "" {
		cli, err := s.client.GetClient(svr.ClientID)
		if err != nil {
			return nil, err
		}
		svr.Client = cli
	}

	return svr, nil
}

// DeleteService todo
func (s *Store) DeleteService(name string) error {
	if name == s.conf.APP.Name {
		return exception.NewBadRequest("your can't delete github.com/defineiot/keyauth service ,beacause can't delete myself")
	}

	svr, err := s.service.GetService(name)
	if err != nil {
		return err
	}

	if svr.ClientID != "" {
		if err := s.client.DeleteClient(svr.ClientID); err != nil {
			return err
		}
	}

	return s.service.DeleteService(name)
}

// RegistryServiceFeatures todo
func (s *Store) RegistryServiceFeatures(name string, features ...service.Feature) error {
	return s.service.RegistryServiceFeatures(name, features...)
}

// ListServiceFeatures todo
func (s *Store) ListServiceFeatures(name string) ([]*service.Feature, error) {
	return s.service.ListServiceFeatures(name)
}

// ListDomainFeatures todo
func (s *Store) ListDomainFeatures() ([]*service.Feature, error) {
	return s.service.ListDomainFeatures()
}

// ListMemberFeatures todo
func (s *Store) ListMemberFeatures() ([]*service.Feature, error) {
	return s.service.ListMemberFeatures()
}

// CheckServiceHasFeature todo
func (s *Store) CheckServiceHasFeature(sn, fn string) (bool, error) {
	return s.service.CheckServiceHasFeature(sn, fn)
}

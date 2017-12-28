package mock

import "openauth/store/domain"

// DomainStore mock
type DomainStore struct {
	CreateDomainFn                  func(name, description, displayName string, enabled bool) (*domain.Domain, error)
	CreateDomainInvoked             bool
	GetDomainFn                     func(domainID string) (*domain.Domain, error)
	GetDomainInvoked                bool
	GetDomainByNameFn               func(name string) (*domain.Domain, error)
	GetDomainByNameInvoked          bool
	ListDomainFn                    func(pageNumber, pageSize int64) ([]*domain.Domain, int64, error)
	ListDomainInvoked               bool
	UpdateDomainFn                  func(id, name, description string) (*domain.Domain, error)
	UpdateDomainInvoked             bool
	DeleteDomainFn                  func(id string) error
	DeleteDomainInvoked             bool
	CheckDomainIsExistByIDFn        func(domainID string) (bool, error)
	CheckDomainIsExistByIDInvoked   bool
	CheckDomainIsExistByNameFn      func(domainName string) (bool, error)
	CheckDomainIsExistByNameInvoked bool
}

// Close mock
func (s *DomainStore) Close() error {
	return nil
}

// CreateDomain mock
func (s *DomainStore) CreateDomain(name, description, displayName string, enabled bool) (*domain.Domain, error) {
	s.CreateDomainInvoked = true
	d, err := s.CreateDomainFn(name, description, displayName, enabled)
	return d, err
}

// GetDomain mock
func (s *DomainStore) GetDomain(domainID string) (*domain.Domain, error) {
	s.GetDomainInvoked = true
	d, err := s.GetDomainFn(domainID)
	return d, err
}

// GetDomainByName mock
func (s *DomainStore) GetDomainByName(name string) (*domain.Domain, error) {
	s.GetDomainByNameInvoked = true
	d, err := s.GetDomainByNameFn(name)
	return d, err
}

// ListDomain mock
func (s *DomainStore) ListDomain(pageNumber, pageSize int64) ([]*domain.Domain, int64, error) {
	s.ListDomainInvoked = true
	domains, totals, err := s.ListDomainFn(pageNumber, pageSize)
	return domains, totals, err
}

// UpdateDomain mock
func (s *DomainStore) UpdateDomain(id, name, description string) (*domain.Domain, error) {
	s.UpdateDomainInvoked = true
	d, err := s.UpdateDomainFn(id, name, description)
	return d, err
}

// DeleteDomain mock
func (s *DomainStore) DeleteDomain(id string) error {
	s.DeleteDomainInvoked = true
	err := s.DeleteDomainFn(id)
	return err
}

// CheckDomainIsExistByID mock
func (s *DomainStore) CheckDomainIsExistByID(domainID string) (bool, error) {
	s.CheckDomainIsExistByIDInvoked = true
	ok, err := s.CheckDomainIsExistByIDFn(domainID)
	return ok, err
}

// CheckDomainIsExistByName mock
func (s *DomainStore) CheckDomainIsExistByName(domainName string) (bool, error) {
	s.CheckDomainIsExistByNameInvoked = true
	ok, err := s.CheckDomainIsExistByNameFn(domainName)
	return ok, err
}

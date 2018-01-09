package mock

import "openauth/store/service"

// ServiceStore mock
type ServiceStore struct {
	SaveServiceFn          func(name, description string) (*service.Service, error)
	SaveServiceInvoked     bool
	DeleteServiceFn        func(sid string) error
	DeleteServiceInvoked   bool
	FindAllServiceFn       func() ([]*service.Service, error)
	FindAllServiceInvoked  bool
	FindServiceByIDFn      func(sid string) (*service.Service, error)
	FindServiceByIDInvoked bool
}

// Close mock
func (s *ServiceStore) Close() error {
	return nil
}

// SaveService mock
func (s *ServiceStore) SaveService(name, description string) (*service.Service, error) {
	s.SaveServiceInvoked = true
	svr, err := s.SaveServiceFn(name, description)
	return svr, err
}

// DeleteService mock
func (s *ServiceStore) DeleteService(sid string) error {
	s.DeleteServiceInvoked = true
	return s.DeleteServiceFn(sid)
}

// FindAllService mock
func (s *ServiceStore) FindAllService() ([]*service.Service, error) {
	s.FindAllServiceInvoked = true
	svrs, err := s.FindAllServiceFn()
	return svrs, err
}

// FindServiceByID mock
func (s *ServiceStore) FindServiceByID(sid string) (*service.Service, error) {
	s.FindServiceByIDInvoked = true
	svr, err := s.FindServiceByIDFn(sid)
	return svr, err
}

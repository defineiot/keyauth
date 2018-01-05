package mysql

import (
	"time"

	"github.com/satori/go.uuid"

	"openauth/store/service"
)

func (s *store) SaveService(name, description string) (*service.Service, error) {
	svr := new(service.Service)
	svr.ID = uuid.NewV4().String()
	svr.Name = name
	svr.Description = description
	svr.CreateAt = time.Now().Unix()
	svr.Enabled = true

	return nil, nil
}

func (s *store) DeleteService(sid string) error {
	return nil
}

func (s *store) FindAllService() ([]*service.Service, error) {
	return nil, nil
}

func (s *store) FindServiceByID() (*service.Service, error) {
	return nil, nil
}

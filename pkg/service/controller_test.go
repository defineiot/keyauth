package service_test

import (
	"fmt"
	"openauth/pkg/service"
	svrstr "openauth/store/service"
	"openauth/store/service/mock"

	"github.com/sirupsen/logrus"
)

func NewServiceController() *service.Controller {
	log := logrus.New()
	ss := new(mock.ServiceStore)

	mockSvr := svrstr.Service{
		ID:   "validated-sid",
		Name: "validated-name",
	}
	ss.SaveServiceFn = func(name, description string) (*svrstr.Service, error) {
		svr := svrstr.Service{
			Name:        name,
			Description: description,
		}

		return &svr, nil
	}
	ss.DeleteServiceFn = func(sid string) error {
		if sid == "validated-sid" {
			return nil
		}
		return fmt.Errorf("%s not found", sid)
	}
	ss.FindAllServiceFn = func() ([]*svrstr.Service, error) {
		return []*svrstr.Service{&mockSvr}, nil
	}

	ss.FindServiceByIDFn = func(sid string) (*svrstr.Service, error) {
		if sid != "validated-sid" {
			return nil, fmt.Errorf("%s not found", sid)
		}
		return &mockSvr, nil
	}

	return service.NewController(ss, log)
}

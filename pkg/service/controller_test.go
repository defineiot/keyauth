package service_test

import (
	"openauth/pkg/service"
	svrstr "openauth/store/service"
	"openauth/store/service/mock"

	"github.com/sirupsen/logrus"
)

func NewServiceController() *service.Controller {
	log := logrus.New()
	ss := new(mock.ServiceStore)

	ss.SaveServiceFn = func(name, description string) (*svrstr.Service, error) {
		svr := svrstr.Service{
			Name:        name,
			Description: description,
		}

		return &svr, nil
	}

	return service.NewController(ss, log)
}

package domain_test

import (
	"errors"

	"github.com/sirupsen/logrus"

	"openauth/pkg/domain"
	domainstr "openauth/store/domain"
	"openauth/store/domain/mock"
)

func NewDomainController() *domain.Controller {
	log := logrus.New()
	ds := new(mock.DomainStore)
	ds.CreateDomainFn = func(name, description, displayName string, enabled bool) (*domainstr.Domain, error) {
		dom := domainstr.Domain{
			Name:        name,
			Description: description,
			DisplayName: displayName,
			Enabled:     enabled,
		}

		return &dom, nil
	}

	dom := domainstr.Domain{
		Name:    "unit-test-01",
		Enabled: true,
		ID:      "unit-test-id",
	}
	ds.ListDomainFn = func(pageNumber, pageSize int64) ([]*domainstr.Domain, int64, error) {
		return []*domainstr.Domain{&dom}, 1, nil
	}

	ds.GetDomainFn = func(domainID string) (*domainstr.Domain, error) {
		if domainID == "unit-test-id" {
			return &dom, nil
		}

		return nil, errors.New("unit-test-id not found")
	}
	ds.DeleteDomainFn = func(id string) error {
		if id == "unit-test-id" {
			return nil
		}

		return errors.New("unit-test-id not found")
	}

	return domain.NewController(ds, log)
}

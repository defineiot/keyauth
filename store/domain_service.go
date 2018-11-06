package store

import (
	"errors"
	"strings"

	"github.com/defineiot/keyauth/dao/domain"
	"github.com/defineiot/keyauth/internal/exception"
)

// CreateDomain use to create an domain
func (s *Store) CreateDomain(name, description, displayName string, enabled bool) (*domain.Domain, error) {
	return s.domain.CreateDomain(name, description, displayName, enabled)
}

// GetDomain use to get domain by id or name
func (s *Store) GetDomain(bywhat string, valule string) (*domain.Domain, error) {
	var err error

	dom := new(domain.Domain)
	cacheKey := "domain_" + valule

	if s.isCache {
		if s.cache.Get(cacheKey, dom) {
			s.log.Debugf("get domain from cache key: %s", cacheKey)
			return dom, nil
		}
		s.log.Debugf("get domain from cache failed, key: %s", cacheKey)
	}

	switch strings.ToLower(bywhat) {
	case "id":
		dom, err = s.domain.GetDomainByID(valule)
	case "name":
		dom, err = s.domain.GetDomainByName(valule)
	default:
		return nil, errors.New("only support ID and Name to get Domain")
	}

	if err != nil {
		return nil, err
	}
	if dom == nil {
		return nil, exception.NewBadRequest("domain: %s not fond", valule)
	}

	if s.isCache {
		if !s.cache.Set(cacheKey, dom, s.ttl) {
			s.log.Debugf("set domain cache failed, key: %s", cacheKey)
		}
		s.log.Debugf("set domain cache ok, key: %s", cacheKey)
	}

	return dom, nil
}

// ListDomain all domain
func (s *Store) ListDomain(pageNumber, pageSize int64) (domains []*domain.Domain, totalPage int64, err error) {
	return s.domain.ListDomain(pageNumber, pageSize)
}

// DeleteDomain an exist domain
func (s *Store) DeleteDomain(bywhat string, value string) error {
	var err error

	cacheKey := "domain_" + value

	switch strings.ToLower(bywhat) {
	case "id":
		err = s.domain.DeleteDomainByID(value)
	case "name":
		err = s.domain.DeleteDomainByName(value)
	default:
		return errors.New("only support ID and Name to get Domain")
	}

	if err != nil {
		return err
	}

	if s.isCache {
		if !s.cache.Delete(cacheKey) {
			s.log.Debugf("delete domain from cache failed, key: %s", cacheKey)
		}
		s.log.Debugf("delete domain from cache suucess, key: %s", cacheKey)
	}

	return nil
}

// CheckDomainExistByName check domain exist
func (s *Store) CheckDomainExistByName(name string) (bool, error) {
	return s.domain.CheckDomainIsExistByName(name)
}

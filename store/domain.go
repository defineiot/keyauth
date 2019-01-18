package store

import (
	"errors"
	"strings"

	"github.com/defineiot/keyauth/dao/domain"
	"github.com/defineiot/keyauth/internal/exception"
)

// CreateDomain use to create an domain
func (s *Store) CreateDomain(d *domain.Domain) error {
	return s.dao.Domain.CreateDomain(d)
}

// GetDomain use to get domain by id or name
func (s *Store) GetDomain(bywhat string, valule string) (*domain.Domain, error) {
	var err error

	dom := new(domain.Domain)
	cacheKey := "domain_" + valule

	if s.isCache {
		if s.cache.Get(cacheKey, dom) {
			s.log.Debug("get domain from cache key: %s", cacheKey)
			return dom, nil
		}
		s.log.Debug("get domain from cache failed, key: %s", cacheKey)
	}

	switch strings.ToLower(bywhat) {
	case "id":
		dom, err = s.dao.Domain.GetDomainByID(valule)
	case "name":
		dom, err = s.dao.Domain.GetDomainByName(valule)
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
			s.log.Debug("set domain cache failed, key: %s", cacheKey)
		}
		s.log.Debug("set domain cache ok, key: %s", cacheKey)
	}

	return dom, nil
}

// ListDomain all domain
func (s *Store) ListDomain(pageNumber, pageSize int64) (domains []*domain.Domain, totalPage int64, err error) {
	return s.dao.Domain.ListDomain(pageNumber, pageSize)
}

// DeleteDomain an exist domain
func (s *Store) DeleteDomain(bywhat string, value string) error {
	var err error

	cacheKey := "domain_" + value

	switch strings.ToLower(bywhat) {
	case "id":
		err = s.dao.Domain.DeleteDomainByID(value)
	case "name":
		err = s.dao.Domain.DeleteDomainByName(value)
	default:
		return errors.New("only support ID and Name to get Domain")
	}

	if err != nil {
		return err
	}

	if s.isCache {
		if !s.cache.Delete(cacheKey) {
			s.log.Debug("delete domain from cache failed, key: %s", cacheKey)
		}
		s.log.Debug("delete domain from cache suucess, key: %s", cacheKey)
	}

	return nil
}

// CheckDomainExistByName check domain exist
func (s *Store) CheckDomainExistByName(name string) (bool, error) {
	return s.dao.Domain.CheckDomainIsExistByName(name)
}

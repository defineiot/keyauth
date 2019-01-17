package store

import (
	"time"

	"github.com/defineiot/keyauth/dao"

	"github.com/defineiot/keyauth/internal/cache"
	"github.com/defineiot/keyauth/internal/conf"
	"github.com/defineiot/keyauth/internal/logger"

	mysqlDom "github.com/defineiot/keyauth/dao/domain/mysql"
	mysqlPro "github.com/defineiot/keyauth/dao/project/mysql"
)

// Store is DAO
type Store struct {
	*dao.Dao
	log     logger.Logger
	cache   cache.Cache
	ttl     time.Duration
	conf    *conf.Config
	isCache bool
}

// NewStore store engine
func NewStore(conf *conf.Config) (*Store, error) {
	db, err := conf.GetDBConn()
	if err != nil {
		return nil, err
	}

	store := new(Store)
	dom, err := mysqlDom.NewDomainStore(db)
	if err != nil {
		return nil, err
	}
	pro, err := mysqlPro.NewProjectStore(db)
	if err != nil {
		return nil, err
	}

	return nil, &store
}

// SetCache set cache instance
func (s *Store) SetCache(cache cache.Cache, ttl time.Duration) {
	s.cache = cache
	s.ttl = ttl
	s.isCache = true
}

// DisableCache cancel cache
func (s *Store) DisableCache() {
	s.isCache = false
}

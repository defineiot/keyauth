package store

import (
	"time"

	"github.com/defineiot/keyauth/dao/application"
	"github.com/defineiot/keyauth/dao/client"
	"github.com/defineiot/keyauth/dao/domain"
	"github.com/defineiot/keyauth/dao/project"
	"github.com/defineiot/keyauth/dao/role"
	"github.com/defineiot/keyauth/dao/service"
	"github.com/defineiot/keyauth/dao/token"
	"github.com/defineiot/keyauth/dao/user"
	"github.com/defineiot/keyauth/internal/cache"
	"github.com/defineiot/keyauth/internal/conf"
	"github.com/defineiot/keyauth/internal/log"

	mysqlDom "github.com/defineiot/keyauth/dao/domain/mysql"
	mysqlPro "github.com/defineiot/keyauth/dao/project/mysql"
)

// Store is DAO
type Store struct {
	domain  domain.Store
	project project.Store
	user    user.Store
	app     application.Store
	client  client.Store
	token   token.Store
	service service.Store
	role    role.Store

	log     log.IOTAuthLogger
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

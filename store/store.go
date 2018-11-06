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

// Options new store options
type Options struct {
	Domain  domain.Store
	Project project.Store
	User    user.Store
	App     application.Store
	Client  client.Store
	Token   token.Store
	Role    role.Store
	Service service.Store
	Log     log.IOTAuthLogger
	Conf    *conf.Config
}

// NewStore store engine
func NewStore(opts *Options) *Store {
	store := Store{domain: opts.Domain, log: opts.Log, project: opts.Project, user: opts.User, app: opts.App, client: opts.Client, token: opts.Token, conf: opts.Conf, service: opts.Service, role: opts.Role}

	return &store
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

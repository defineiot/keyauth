package store

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"time"

	"github.com/defineiot/keyauth/dao"
	_ "github.com/defineiot/keyauth/dao/all"

	"github.com/defineiot/keyauth/internal/cache"
	"github.com/defineiot/keyauth/internal/conf"
	"github.com/defineiot/keyauth/internal/logger"
)

// Store is DAO
type Store struct {
	dao     *dao.Dao
	log     logger.Logger
	cache   cache.Cache
	ttl     time.Duration
	conf    *conf.Config
	isCache bool

	tokenCachePrefix string
}

// NewStore store engine
func NewStore(conf *conf.Config) (*Store, error) {
	db, err := conf.GetDBConn()
	if err != nil {
		panic(err)
	}

	log, err := conf.GetLogger()
	if err != nil {
		panic(err)
	}

	opt := &dao.Options{DB: db, LOG: log}
	defaultDao, err := dao.Init(opt)
	if err != nil {
		return nil, err
	}

	s := new(Store)
	s.conf = conf
	s.log = log
	s.dao = defaultDao

	s.tokenCachePrefix = "token_"

	return s, nil
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

func (s *Store) hmacHash(msg string) string {
	mac := hmac.New(sha256.New, []byte(s.conf.APP.Key))
	io.WriteString(mac, msg)

	return fmt.Sprintf("%x", mac.Sum(nil))
}

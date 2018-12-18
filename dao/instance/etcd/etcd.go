package etcd

import (
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"

	"github.com/defineiot/keyauth/dao/instance"
)

// NewInstanceStore todo
func NewInstanceStore(endpoints []string, registryPrefix string) (instance.Store, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Duration(5) * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("connect etcd error, %s", err)
	}

	return &store{
		client:         cli,
		prefix:         registryPrefix,
		requestTimeout: time.Duration(5) * time.Second,
	}, nil
}

// DomainManager is use mongodb as storage
type store struct {
	requestTimeout time.Duration
	client         *clientv3.Client
	prefix         string
}

// Close closes the database, releasing any open resources.
func (s *store) Close() error {
	return s.client.Close()
}

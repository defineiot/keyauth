package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"

	"github.com/defineiot/keyauth/api/global"
	"github.com/defineiot/keyauth/api/register"
)

type etcd struct {
	leaseID        clientv3.LeaseID
	client         *clientv3.Client
	requestTimeout time.Duration

	isStopped    bool
	instanceKey  string
	stopInstance chan struct{}
}

// NewEtcdRegister 初始化一个基于etcd的实例注册器
func NewEtcdRegister(endpoints []string) (register.Register, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Duration(5) * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("connect etcd error, %s", err)
	}

	etcdR := new(etcd)
	etcdR.client = client
	etcdR.stopInstance = make(chan struct{}, 1)
	etcdR.requestTimeout = time.Duration(5) * time.Second

	return etcdR, nil
}

// Register use to registe serice endpoint to etcd. when etcd is down,
// register can retry to registe util the etcd up.
//
// name is service name, use to discovery service address, eg. keyauth
// host and port is service endpoint, eg. 127.0.0.0:50000
// target is etcd addr, eg. 127.0.0.0:2379
// interval is service refresh interval, eg. 10s
// ttl is service ttl, eg. 15
func (e *etcd) Registe(service *register.ServiceInstance) error {
	if err := service.Validate(); err != nil {
		return err
	}

	sjson, err := json.Marshal(service)
	if err != nil {
		global.Log.Error("marshal service object to json error, %s", err)
	}
	serviceValue := string(sjson)

	// minimum lease TTL is ttl-second
	ctxGrant, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	// get lease object
	resp, err := e.client.Lease.Grant(ctxGrant, int64(service.TTL))
	if err != nil {
		global.Log.Warn("service '%s' connect to etcd3 failed: %s, retry after %s second.", service.Name, err.Error(), service.Interval)
		return err
	}
	e.leaseID = resp.ID

	// lease keep alive
	keepevent, err := e.client.Lease.KeepAlive(context.TODO(), resp.ID)
	if err != nil {
		global.Log.Error("lease keep alive error, %s", err)
	}

	go func() {
		for {
			select {
			case _, ok := <-keepevent:
				if !ok {
					global.Log.Error("keep alive channel closed")
					return
				}
			}
		}
	}()

	// service not registe, set it
	serviceKey := fmt.Sprintf("%s/%s/%s", service.Prefix, service.Type, service.Name)
	if _, err := e.client.Put(context.Background(), serviceKey, serviceValue, clientv3.WithLease(resp.ID)); err != nil {
		global.Log.Warn("registe service '%s' with ttl to etcd3 failed: %s", service.Name, err.Error())
	}
	global.Log.Info("registe service %s success, %s -> %s", service.Name, serviceKey, serviceValue)
	e.instanceKey = serviceKey

	return nil
}

// UnRegiste delete registered service from etcd, if etcd is down
// unregister while timeout.
func (e *etcd) UnRegiste() error {
	if e.isStopped {
		return errors.New("the instance has unregisted")
	}

	// delete instance key
	e.stopInstance <- struct{}{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if resp, err := e.client.Delete(ctx, e.instanceKey); err != nil {
		global.Log.Warn("unregister '%s' failed: connect to etcd server timeout, %s", e.instanceKey, err.Error())
	} else {
		if resp.Deleted == 0 {
			global.Log.Info("unregister '%s' failed, the key not exist", e.instanceKey)
		} else {
			global.Log.Info("unregister '%s' ok", e.instanceKey)
		}
	}

	// revoke lease
	_, err := e.client.Lease.Revoke(context.TODO(), e.leaseID)
	if err != nil {
		global.Log.Warn("revoke lease error, %s", err)
		return err
	}
	e.isStopped = true

	return nil
}

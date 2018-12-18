package etcd_test

import (
	"time"

	"github.com/defineiot/keyauth/dao/instance"
	"github.com/defineiot/keyauth/dao/instance/etcd"
	"github.com/defineiot/keyauth/dao/service"
	"github.com/defineiot/keyauth/internal/conf/mock"
)

func newTestStore() instance.Store {
	conf := mock.NewConfig()

	store, err := etcd.NewInstanceStore(conf.Etcd.Endpoints, conf.Etcd.InstanceRegistryPrefix)
	if err != nil {
		panic(err)
	}

	return store
}

type instanceSuit struct {
	store instance.Store
	ins   *instance.Instance
}

func (s *instanceSuit) TearDown() {
	s.store.Close()
}

func (s *instanceSuit) SetUp() {
	ins := &instance.Instance{
		ID:        "resource_service_aabbxxcc",
		Name:      "r01",
		Address:   "127.0.0.1:8080",
		Version:   "1.0.0.1",
		GitBranch: "master",
		GitCommit: "b5fa6ac8fa056aa18ed2fc6a0eb12497e7db1be1",
		BuildEnv:  "go 1.10.4",
		BuildAt:   "2018-12-31",
		ServiceID: "service-00xxbbcc",
		OnlineAt:  time.Now().Unix(),
	}

	f1 := &service.Feature{
		ID:          "b5fa6ac8fa056aa18ed2fc6a0eb12497e7db1be1",
		Name:        "CreateDevice",
		Method:      "POST",
		Endpoint:    "/resurce_service/devices/",
		Description: "unit-test-for-feature",
	}
	f2 := &service.Feature{
		ID:          "718a4ed919c0fc59a5157b02e594ff2411c041c0",
		Name:        "GetDevice",
		Method:      "GET",
		Endpoint:    "/resurce_service/devices/",
		Description: "unit-test-for-feature",
	}
	ins.Features = []*service.Feature{f1, f2}

	s.ins = ins
	s.store = newTestStore()

}

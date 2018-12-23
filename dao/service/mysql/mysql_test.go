package mysql_test

import (
	"github.com/defineiot/keyauth/dao/service"
	"github.com/defineiot/keyauth/dao/service/mysql"
	"github.com/defineiot/keyauth/internal/conf/mock"
)

func newTestStore() service.Store {
	conf := mock.NewConfig()
	db, err := conf.GetDBConn()
	if err != nil {
		panic(err)
	}

	log, err := conf.GetLogger()
	if err != nil {
		panic(err)
	}

	store, err := mysql.NewServiceStore(db, log)
	if err != nil {
		panic(err)
	}

	return store
}

type serviceSuit struct {
	roleID string

	store    service.Store
	svr      *service.Service
	features []*service.Feature
}

func (s *serviceSuit) TearDown() {
	s.store.Close()
}

func (s *serviceSuit) SetUp() {
	s.roleID = "unit-test-role-id"

	s.svr = &service.Service{
		Name:        "unit-test-service-name01",
		Description: "just for unit test",
		Type:        service.Public,
	}

	f1 := &service.Feature{
		Name:         "feature01",
		Tag:          "POST",
		HTTPEndpoint: "/features/01",
		Description:  "only for unit test",
	}
	f2 := &service.Feature{
		Name:         "feature02",
		Tag:          "POST",
		HTTPEndpoint: "/features/02",
		Description:  "only for unit test",
	}
	s.features = []*service.Feature{f1, f2}

	s.store = newTestStore()

}

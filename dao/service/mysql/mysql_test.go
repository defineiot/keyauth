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
	store service.Store
	svr   *service.Service
}

func (s *serviceSuit) TearDown() {
	s.store.Close()
}

func (s *serviceSuit) SetUp() {
	s.svr = &service.Service{
		Name:        "unit-test-service-name01",
		Description: "just for unit test",
		Type:        service.Public,
	}

	s.store = newTestStore()

}

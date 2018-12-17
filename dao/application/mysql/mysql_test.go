package mysql_test

import (
	"github.com/defineiot/keyauth/dao/application"
	"github.com/defineiot/keyauth/dao/application/mysql"
	"github.com/defineiot/keyauth/internal/conf/mock"
)

func newTestStore() application.Store {
	conf := mock.NewConfig()
	db, err := conf.GetDBConn()
	if err != nil {
		panic(err)
	}

	store, err := mysql.NewAppStore(db)
	if err != nil {
		panic(err)
	}

	return store
}

type applicationSuit struct {
	store  application.Store
	app    *application.Application
	userID string
}

func (s *applicationSuit) TearDown() {
	s.store.Close()
}

func (s *applicationSuit) SetUp() {
	s.store = newTestStore()
	s.userID = "unit-test-01"
	s.app = &application.Application{
		Name:   "application01",
		UserID: s.userID,
	}
}

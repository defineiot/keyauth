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

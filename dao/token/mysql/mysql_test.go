package mysql_test

import (
	"github.com/defineiot/keyauth/internal/conf/mock"
	"github.com/defineiot/keyauth/store/token"
	"github.com/defineiot/keyauth/store/token/mysql"
)

func newTestStore() token.Store {
	conf := mock.NewConfig()
	db, err := conf.GetDBConn()
	if err != nil {
		panic(err)
	}

	store, err := mysql.NewTokenStore(db)
	if err != nil {
		panic(err)
	}

	return store
}

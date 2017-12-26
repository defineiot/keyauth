package mysql_test

import (
	"openauth/api/config/mock"
	"openauth/store/token"
	"openauth/store/token/mysql"
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

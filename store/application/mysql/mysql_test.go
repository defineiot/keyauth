package mysql_test

import (
	"openauth/api/config/mock"
	"openauth/store/application"
	"openauth/store/application/mysql"
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

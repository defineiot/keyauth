package mysql_test

import (
	"openauth/api/config/mock"
	"openauth/store/service"
	"openauth/store/service/mysql"
)

func newTestStore() service.Store {
	conf := mock.NewConfig()
	db, err := conf.GetDBConn()
	if err != nil {
		panic(err)
	}

	store, err := mysql.NewServiceStore(db)
	if err != nil {
		panic(err)
	}

	return store
}

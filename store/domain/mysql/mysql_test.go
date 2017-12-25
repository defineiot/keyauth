package mysql_test

import (
	"openauth/api/config/mock"
	"openauth/store/domain"
	"openauth/store/domain/mysql"
)

func newTestStore() domain.Store {
	conf := mock.NewConfig()
	db, err := conf.GetDBConn()
	if err != nil {
		panic(err)
	}

	store, err := mysql.NewDomainStore(db)
	if err != nil {
		panic(err)
	}

	return store
}

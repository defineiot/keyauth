package mysql_test

import (
	"openauth/api/config/mock"
	"openauth/store/project"
	"openauth/store/project/mysql"
)

func newTestStore() project.Store {
	conf := mock.NewConfig()
	db, err := conf.GetDBConn()
	if err != nil {
		panic(err)
	}

	store, err := mysql.NewProjectStore(db)
	if err != nil {
		panic(err)
	}

	return store
}

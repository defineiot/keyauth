package mysql_test

import (
	"github.com/defineiot/keyauth/dao/project"
	"github.com/defineiot/keyauth/dao/project/mysql"
	"github.com/defineiot/keyauth/internal/conf/mock"
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

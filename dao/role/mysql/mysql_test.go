package mysql_test

import (
	"github.com/defineiot/keyauth/internal/conf/mock"
	"github.com/defineiot/keyauth/store/role"
	"github.com/defineiot/keyauth/store/role/mysql"
)

func newTestStore() role.Store {
	conf := mock.NewConfig()
	db, err := conf.GetDBConn()
	if err != nil {
		panic(err)
	}

	log, err := conf.GetLogger()
	if err != nil {
		panic(err)
	}

	store, err := mysql.NewRoleStore(db, log)
	if err != nil {
		panic(err)
	}

	return store
}

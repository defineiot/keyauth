package mysql_test

import (
	"github.com/defineiot/keyauth/dao/domain"
	"github.com/defineiot/keyauth/dao/domain/mysql"
	"github.com/defineiot/keyauth/internal/conf/mock"
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

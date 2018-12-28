package mysql_test

import (
	"github.com/sirupsen/logrus"

	"github.com/defineiot/keyauth/dao/user"
	"github.com/defineiot/keyauth/dao/user/mysql"
	"github.com/defineiot/keyauth/internal/conf/mock"
)

func newTestStore() user.Store {
	conf := mock.NewConfig()
	db, err := conf.GetDBConn()
	if err != nil {
		panic(err)
	}

	log := logrus.New()
	store, err := mysql.NewUserStore(db, "default", log)
	if err != nil {
		panic(err)
	}

	return store
}

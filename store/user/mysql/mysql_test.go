package mysql_test

import (
	"github.com/sirupsen/logrus"

	"openauth/api/config/mock"
	"openauth/store/user"
	"openauth/store/user/mysql"
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

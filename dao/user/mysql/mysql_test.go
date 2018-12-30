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

type userSuit struct {
	u     *user.User
	store user.Store
}

func (s *userSuit) TearDown() {
	s.store.Close()
}

func (s *userSuit) SetUp() {
	s.u = &user.User{
		Account:           "unit-test-for-user01",
		Mobile:            "18108054577",
		Email:             "18108054577@163.com",
		Phone:             "028-1111111",
		Address:           "家庭住址",
		RealName:          "单元测试",
		NickName:          "陈独秀",
		Gender:            "M",
		Avatar:            "www.google.com",
		Language:          "zh_CN",
		City:              "成都",
		Province:          "四川",
		ExpiresActiveDays: 90,
		Password:          &user.Password{Password: "123456", ExpireAt: 365},
		DomainID:          "unit-test-for-domain",
		DepartmentID:      "unit-test-for-department",
	}

	s.store = newTestStore()

}

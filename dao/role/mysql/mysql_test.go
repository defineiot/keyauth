package mysql_test

import (
	"github.com/defineiot/keyauth/dao"
	"github.com/defineiot/keyauth/dao/models"
	"github.com/defineiot/keyauth/dao/role"
	"github.com/defineiot/keyauth/dao/role/mysql"
	"github.com/defineiot/keyauth/internal/conf/mock"
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

	opt := &dao.Options{DB: db, LOG: log}
	store, err := mysql.NewRoleStore(opt)
	if err != nil {
		panic(err)
	}

	return store
}

type roleSuit struct {
	store role.Store
	r     *models.Role
	name  string
}

func (s *roleSuit) TearDown() {
	s.store.Close()
}

func (s *roleSuit) SetUp() {
	s.name = "role-unit-test-01"

	s.r = &models.Role{
		Name:        s.name,
		Description: "unit-test",
	}

	s.store = newTestStore()

}

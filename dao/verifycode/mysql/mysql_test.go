package mysql_test

import (
	"github.com/defineiot/keyauth/dao"
	"github.com/defineiot/keyauth/dao/models"
	"github.com/defineiot/keyauth/dao/verifycode"
	"github.com/defineiot/keyauth/dao/verifycode/mysql"
	"github.com/defineiot/keyauth/internal/conf/mock"
)

func newTestStore() verifycode.Store {
	conf := mock.NewConfig()
	db, err := conf.GetDBConn()
	if err != nil {
		panic(err)
	}

	opt := &dao.Options{DB: db}
	store, err := mysql.NewVerifyCodeStore(opt)
	if err != nil {
		panic(err)
	}

	return store
}

type verifyCodeSuit struct {
	store  verifycode.Store
	code   *models.VerifyCode
	target string
}

func (s *verifyCodeSuit) TearDown() {
	s.store.Close()
}

func (s *verifyCodeSuit) SetUp() {
	s.target = "18108053819"

	s.code = &models.VerifyCode{
		Purpose:    models.RegistryCode,
		SendMode:   models.SendByMobile,
		SendTarget: s.target,
	}

	s.store = newTestStore()

}

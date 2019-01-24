package dao_test

import (
	"fmt"
	"testing"

	_ "github.com/defineiot/keyauth/dao/all"

	"github.com/defineiot/keyauth/dao"
	"github.com/defineiot/keyauth/dao/token"
	"github.com/defineiot/keyauth/internal/conf/mock"
	"github.com/stretchr/testify/require"
)

func newTestDao() *dao.Dao {
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
	defaultDao, err := dao.Init(opt)
	if err != nil {
		panic(err)
	}

	return defaultDao
}

type daoSuit struct {
	dao *dao.Dao

	tk *token.Token
}

func (s *daoSuit) TearDown() {
}

func (s *daoSuit) SetUp() {
	s.tk = &token.Token{
		AccessToken:   token.MakeBearerToken(24),
		RefreshToken:  token.MakeBearerToken(32),
		TokenType:     token.Bearer,
		GrantType:     token.PASSWORD,
		UserID:        "dao-unit-test",
		DomainID:      "dao-unit-test",
		ApplicationID: "dao-unit-test",
		ExpiresIn:     3600,
	}

	s.dao = newTestDao()
}

func TestDAOtSuit(t *testing.T) {
	suit := new(daoSuit)
	suit.SetUp()
	defer suit.TearDown()

	t.Run("CreateTokenOK", testCreateTokenOK(suit))
	t.Run("DeleteTokenByAccessTOkenOK", testDeleteTokenByAccessTokenOK(suit))
}

func testCreateTokenOK(s *daoSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		fmt.Println(s.dao)
		err := s.dao.Token.SaveToken(s.tk)
		should.NoError(err)

		t.Logf("save token success: %s", s.tk)
	}
}

func testDeleteTokenByAccessTokenOK(s *daoSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		err := s.dao.Token.DeleteToken(s.tk.AccessToken)
		should.NoError(err)

		t.Logf("delete token by access token success: %s", s.tk)
	}
}

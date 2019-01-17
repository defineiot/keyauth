package mysql_test

import (
	"github.com/defineiot/keyauth/dao"
	"github.com/defineiot/keyauth/dao/token"
	"github.com/defineiot/keyauth/dao/token/mysql"
	"github.com/defineiot/keyauth/internal/conf/mock"
)

func newTestStore() token.Store {
	conf := mock.NewConfig()
	db, err := conf.GetDBConn()
	if err != nil {
		panic(err)
	}

	opt := &dao.Options{DB: db}
	store, err := mysql.NewTokenStore(opt)
	if err != nil {
		panic(err)
	}

	return store
}

type tokenSuit struct {
	uid string
	did string
	aid string

	t1    *token.Token
	t2    *token.Token
	store token.Store
}

func (s *tokenSuit) TearDown() {
	s.store.Close()
}

func (s *tokenSuit) SetUp() {
	s.uid = "unit-test-for-token-user-id"
	s.did = "unit-test-for-token-domain-id"
	s.aid = "unit-test-for-token-application-id"

	s.t1 = &token.Token{
		AccessToken:   token.MakeBearerToken(24),
		RefreshToken:  token.MakeBearerToken(32),
		TokenType:     token.Bearer,
		GrantType:     token.PASSWORD,
		UserID:        s.uid,
		DomainID:      s.did,
		ApplicationID: s.aid,
		ExpiresIn:     3600,
	}
	s.t2 = &token.Token{
		AccessToken:   token.MakeBearerToken(24),
		RefreshToken:  token.MakeBearerToken(32),
		TokenType:     token.Bearer,
		GrantType:     token.PASSWORD,
		UserID:        s.uid,
		DomainID:      s.did,
		ApplicationID: s.aid,
		ExpiresIn:     3600,
	}

	s.store = newTestStore()

}

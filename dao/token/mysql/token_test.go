package mysql_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTokentSuit(t *testing.T) {
	suit := new(tokenSuit)
	suit.SetUp()
	defer suit.TearDown()

	t.Run("CreateTokenOK", testCreateTokenOK(suit))
	t.Run("GetTokenByAccessTokenOK", testGetTokenByAccessTokenOK(suit))
	t.Run("GetTokenByRefreshTokenOK", testGetTokenByRefreshTokenOK(suit))
	t.Run("DeleteTokenByAccessTOkenOK", testDeleteTokenByAccessTokenOK(suit))
	t.Run("DeleteTokenByRefreshTokenOK", testDeleteTokenByRefreshTokenOK(suit))
}

func testCreateTokenOK(s *tokenSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		err := s.store.SaveToken(s.t1)
		should.NoError(err)
		err = s.store.SaveToken(s.t2)
		should.NoError(err)

		t.Logf("save token success: %s", s.t1)
		t.Logf("save token success: %s", s.t2)
	}
}

func testGetTokenByAccessTokenOK(s *tokenSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		tk, err := s.store.GetToken(s.t1.AccessToken)
		should.NoError(err)

		t.Logf("get token by access token(%s) success: %s", s.t1.AccessToken, tk)
	}
}

func testGetTokenByRefreshTokenOK(s *tokenSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		tk, err := s.store.GetTokenByRefresh(s.t1.RefreshToken)
		should.NoError(err)

		t.Logf("get token by refresh token(%s) success: %s", s.t1.RefreshToken, tk)
	}
}

func testDeleteTokenByAccessTokenOK(s *tokenSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		err := s.store.DeleteToken(s.t1.AccessToken)
		should.NoError(err)

		t.Logf("delete token by access token success: %s", s.t1)
	}
}

func testDeleteTokenByRefreshTokenOK(s *tokenSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		err := s.store.DeleteTokenByRefresh(s.t2.RefreshToken)
		should.NoError(err)

		t.Logf("delete token by refresh token success: %s", s.t2)
	}
}

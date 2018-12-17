package mysql_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/defineiot/keyauth/dao/verifycode"
)

func TestVerifyCodeSuit(t *testing.T) {
	suit := new(verifyCodeSuit)
	suit.SetUp()
	defer suit.TearDown()

	t.Run("CreateVerifyCodeOK", testCreateVerifyCodeOK(suit))
	t.Run("GetVerifyCodeOK", testGetVerifyCodeOK(suit))
	t.Run("DeleteVerifyCodeOK", testDeleteVerifyCodeOK(suit))
}

func testCreateVerifyCodeOK(s *verifyCodeSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		err := s.store.CreateVerifyCode(s.code)
		should.NoError(err)
		should.Equal(6, len(s.code.Code))

		t.Logf("create verifycode(%s) success: %s", s.code.Code, s.code)
	}
}

func testGetVerifyCodeOK(s *verifyCodeSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		code, err := s.store.GetVerifyCode(verifycode.Registry, s.target)
		should.NoError(err)

		t.Logf("get verifycode(%s) success: %s", code.Code, code)
	}
}

func testDeleteVerifyCodeOK(s *verifyCodeSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		err := s.store.DeleteVerifyCode(verifycode.Registry, s.target, s.code.Code)
		should.NoError(err)

		t.Logf("delete verifycode(%s) success: %s", s.code.Code, s.code)
	}
}

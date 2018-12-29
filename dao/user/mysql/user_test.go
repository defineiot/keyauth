package mysql_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProjectSuit(t *testing.T) {
	suit := new(userSuit)
	suit.SetUp()
	defer suit.TearDown()

	t.Run("CreateUserOK", testCreateUserOK(suit))
}

func testCreateUserOK(s *userSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		err := s.store.CreateUser(s.u)
		should.NoError(err)

		t.Logf("create user(%s) success: %s", s.u.Account, s.u)
	}
}

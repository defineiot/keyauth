package mysql_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/defineiot/keyauth/dao/user"
)

func TestProjectSuit(t *testing.T) {
	suit := new(userSuit)
	suit.SetUp()
	defer suit.TearDown()

	t.Run("CreateUserOK", testCreateUserOK(suit))
	t.Run("ListDomainUsersOK", testListDomainUsersOK(suit))
	t.Run("GetUserByIDOK", testGetUserByIDOK(suit))
	t.Run("GetUserByAccountOK", testGetUserByAccountOK(suit))
	t.Run("DeleteUserByID", testDeleteUserByIDOK(suit))
}

func testCreateUserOK(s *userSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		err := s.store.CreateUser(s.u)
		should.NoError(err)

		t.Logf("create user(%s) success: %s", s.u.Account, s.u)
	}
}

func testListDomainUsersOK(s *userSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		users, err := s.store.ListDomainUsers(s.u.DomainID)
		should.NoError(err)

		t.Logf("list domain(%s) users: %s", s.u.DomainID, users)
	}
}

func testGetUserByIDOK(s *userSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		u, err := s.store.GetUser(user.UserID, s.u.ID)
		should.NoError(err)

		t.Logf("get user(%s) by id success: %s", s.u.ID, u)
	}
}

func testGetUserByAccountOK(s *userSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		u, err := s.store.GetUser(user.Account, s.u.Account)
		should.NoError(err)

		t.Logf("get user(%s) by account success: %s", s.u.Account, u)
	}
}

func testDeleteUserByIDOK(s *userSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		err := s.store.DeleteUser(s.u.DomainID, s.u.ID)
		should.NoError(err)

		t.Logf("delete user(%s) by id success: %s", s.u.Account, s.u)
	}
}

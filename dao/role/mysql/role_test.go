package mysql_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProjectSuit(t *testing.T) {
	suit := new(roleSuit)
	suit.SetUp()
	defer suit.TearDown()

	t.Run("CreateRoleOK", testCreateRoleOK(suit))
	t.Run("testGetRoleOK", testGetRoleOK(suit))
	t.Run("testListRoleOK", testListRoleOK(suit))
	t.Run("testDeleteRoleOK", testDeleteRoleOK(suit))
}

func testCreateRoleOK(s *roleSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		err := s.store.CreateRole(s.r)
		should.NoError(err)

		t.Logf("create project(%s) success: %s", s.r.Name, s.r)
	}
}

func testGetRoleOK(s *roleSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		r, err := s.store.GetRole(s.r.ID)
		should.NoError(err)

		t.Logf("get role(%s) success: %s", s.r.ID, r)
	}
}

func testListRoleOK(s *roleSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		roles, err := s.store.ListRole()
		should.NoError(err)

		t.Logf("list roles(%s) success", roles)
	}
}

func testDeleteRoleOK(s *roleSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		err := s.store.DeleteRole(s.r.ID)
		should.NoError(err)

		t.Logf("delete role(%s) success: %s", s.r.ID, s.r)
	}
}

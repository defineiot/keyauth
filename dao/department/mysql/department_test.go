package mysql_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProjectSuit(t *testing.T) {
	suit := new(departmentSuit)
	suit.SetUp()
	defer suit.TearDown()

	t.Run("CreateDepartment", testCreateDepartmentOK(suit))
	t.Run("ListSubDepartments", testListSubDepartmentOK(suit))
	t.Run("HasSubDepeartDelete", testDeleteDepartmentFailed(suit))
	t.Run("DeleteDepartment", testDeleteDepartmentOK(suit))
}

func testCreateDepartmentOK(s *departmentSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		err := s.store.CreateDepartment(s.l1)
		should.NoError(err)
		t.Logf("create department(%s) success: %s", s.l1.Name, s.l1)

		s.l2.ParentID = s.l1.ID
		err = s.store.CreateDepartment(s.l2)
		should.NoError(err)
		t.Logf("create department(%s) success: %s", s.l2.Name, s.l2)

		s.l3.ParentID = s.l2.ID
		err = s.store.CreateDepartment(s.l3)
		should.NoError(err)
		t.Logf("create department(%s) success: %s", s.l3.Name, s.l3)

		s.l4.ParentID = s.l3.ID
		err = s.store.CreateDepartment(s.l4)
		should.NoError(err)
		t.Logf("create department(%s) success: %s", s.l4.Name, s.l4)
	}
}

func testListSubDepartmentOK(s *departmentSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		deps, err := s.store.ListSubDepartments(s.l2.ID)
		should.NoError(err)

		should.Equal(1, len(deps))
		t.Logf("list partment(%s) sub department(%s) success", s.l2.Name, deps)
	}
}

func testDeleteDepartmentFailed(s *departmentSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		err := s.store.DelDepartment(s.l2.ID)
		should.EqualError(err, "the department has 1 sub departments, your should delete them first!")
	}
}

func testDeleteDepartmentOK(s *departmentSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		err := s.store.DelDepartment(s.l4.ID)
		should.NoError(err)

		err = s.store.DelDepartment(s.l3.ID)
		should.NoError(err)
		err = s.store.DelDepartment(s.l2.ID)
		should.NoError(err)

		err = s.store.DelDepartment(s.l1.ID)
		should.NoError(err)
	}
}

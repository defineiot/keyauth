package mysql_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProjectSuit(t *testing.T) {
	suit := new(projectSuit)
	suit.SetUp()
	defer suit.TearDown()

	t.Run("CreateProject", testCreateProjectOK(suit))
	t.Run("GetProectByID", testGetProjectOK(suit))
	t.Run("ListDomainProjects", testListDomainProjectsOK(suit))
	t.Run("AddUserToProject", testAddUserToProject(suit))
	t.Run("ListUserProjectsAddOK", testListUserProjectsAddOK(suit))
	t.Run("RemoveUserFromProject", testRemoveUserFromProject(suit))
	t.Run("ListUserProjectsDelOK", testListUserProjectsRemOK(suit))
	t.Run("DeleteProjectByID", testDeleteProjectOK(suit))
}

func testCreateProjectOK(s *projectSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		err := s.store.CreateProject(s.p)
		should.NoError(err)

		t.Logf("create project(%s) success: %s", s.p.Name, s.p)
	}
}

func testGetProjectOK(s *projectSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		p, err := s.store.GetProjectByID(s.p.ID)
		should.NoError(err)

		t.Logf("get project(%s) success: %s", p.Name, p)
	}
}

func testListDomainProjectsOK(s *projectSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		ps, err := s.store.ListDomainProjects(s.did)
		should.NoError(err)

		t.Logf("list domain projects success: %s", ps)
	}
}

func testAddUserToProject(s *projectSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		err := s.store.AddUsersToProject(s.p.ID, s.uid)
		should.NoError(err)
	}
}

func testListUserProjectsAddOK(s *projectSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		projects, err := s.store.ListUserProjects(s.did, s.uid)
		should.NoError(err)
		should.Equal(len(projects), 1)

		t.Logf("add user(%s) to project ok: %s", s.uid, projects)
	}
}

func testListUserProjectsRemOK(s *projectSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		projects, err := s.store.ListUserProjects(s.did, s.uid)
		should.NoError(err)
		should.Equal(len(projects), 0)

		t.Logf("remove user(%s) to project ok: %s", s.uid, projects)
	}
}

func testRemoveUserFromProject(s *projectSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		err := s.store.RemoveUsersFromProject(s.p.ID, s.uid)
		should.NoError(err)
	}
}

func testDeleteProjectOK(s *projectSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		err := s.store.DeleteProjectByID(s.p.ID)
		should.NoError(err)
	}
}

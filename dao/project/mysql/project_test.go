package mysql_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	projectID string
)

func TestProject(t *testing.T) {
	t.Run("CreateOK", testCreateProjectOK)
	t.Run("CheckOK", testCheckProjectOK)
	t.Run("GetOK", testGetProjectOK)
	t.Run("AddUserOK", testAddUserToProject)
	t.Run("RemoveUserOK", testRemoveUserFromProject)
	t.Run("ListOK", testListProjectOK)
	t.Run("DeleteOK", testListProjectOK)
}

func testCreateProjectOK(t *testing.T) {
	s := newTestStore()
	defer s.Close()
	assert.NotNil(t, s)

	s.DeleteProjectByName("test-project-01", "unit-test-domain-id")
	p1, err := s.CreateProject("unit-test-domain-id", "test-project-01", "", true)
	assert.NoError(t, err)
	assert.NotNil(t, p1)

	assert.Equal(t, "test-project-01", p1.Name)

	projectID = p1.ID
}

func testCheckProjectOK(t *testing.T) {
	s := newTestStore()
	defer s.Close()
	assert.NotNil(t, s)

	ok, err := s.CheckProjectIsExistByID(projectID)
	assert.NoError(t, err)
	assert.Equal(t, true, ok, "the project not create success")
}

func testGetProjectOK(t *testing.T) {
	s := newTestStore()
	defer s.Close()
	assert.NotNil(t, s)

	pGet, err := s.GetProject(projectID)
	assert.NoError(t, err)

	assert.Equal(t, projectID, pGet.ID)
}

func testListProjectOK(t *testing.T) {
	s := newTestStore()
	defer s.Close()

	projects, err := s.ListDomainProjects("unit-test-domain-id")
	assert.NoError(t, err)

	assert.Equal(t, 1, len(projects), "create 1 project in this domain, but not get 1")
}

func testDeleteProjectOK(t *testing.T) {
	s := newTestStore()
	defer s.Close()

	err := s.DeleteProjectByID(projectID)
	assert.NoError(t, err)

}

func testAddUserToProject(t *testing.T) {
	s := newTestStore()
	defer s.Close()
	assert.NotNil(t, s)

	err := s.AddUsersToProject(projectID, "test-user-01", "test-user-02")
	assert.NoError(t, err)

	uids, err := s.ListProjectUsers(projectID)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(uids), "add 2 user, but not equal 2")
}

func testRemoveUserFromProject(t *testing.T) {
	s := newTestStore()
	defer s.Close()
	assert.NotNil(t, s)

	err := s.RemoveUsersFromProject(projectID, "test-user-01", "test-user-02")
	assert.NoError(t, err)

	uids, err := s.ListProjectUsers(projectID)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(uids), "remove 2 user, but user not equal 0")
}

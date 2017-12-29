package project_test

import (
	"testing"
)

func TestCreateProject(t *testing.T) {
	t.Run("OK", testCreateOK)
	t.Run("DomainNotFound", testCreateDomainNotFound)
}

func testCreateOK(t *testing.T) {
	prosvr := NewProjectController()

	_, err := prosvr.CreateProject("validated-domain-id", "validated-name", "")
	if err != nil {
		t.Fatal(err)
	}
}

func testCreateDomainNotFound(t *testing.T) {
	prosvr := NewProjectController()

	_, err := prosvr.CreateProject("invalidated-domain-id", "validated-name", "")
	if err.Error() != "domain invalidated-domain-id not exist" {
		t.Fatal("want domain not found")
	}
}

func TestListDomain(t *testing.T) {
	t.Run("OK", testListDomainOK)
	t.Run("DomainNotFound", testListDomainNotFound)
}

func testListDomainOK(t *testing.T) {
	prosvr := NewProjectController()

	_, err := prosvr.ListProject("validated-domain-id")
	if err != nil {
		t.Fatal(err)
	}
}

func testListDomainNotFound(t *testing.T) {
	prosvr := NewProjectController()

	_, err := prosvr.ListProject("invalidated-domain-id")
	if err.Error() != "domain invalidated-domain-id not exist" {
		t.Fatal("want domain not found")
	}
}

func TestGetProject(t *testing.T) {
	t.Run("OK", testGetOK)
}

func testGetOK(t *testing.T) {
	prosvr := NewProjectController()

	_, err := prosvr.GetProject("validated-project-id")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDestroyProject(t *testing.T) {
	t.Run("OK", testDestroyOK)
}

func testDestroyOK(t *testing.T) {
	prosvr := NewProjectController()

	if err := prosvr.DestroyProject("validated-id"); err != nil {
		t.Fatal(err)
	}
}

func TestAddUsersToProject(t *testing.T) {
	t.Run("OK", testAddOK)
	t.Run("ProjectNotFound", testProjectNotFound)
}

func testAddOK(t *testing.T) {
	prosvr := NewProjectController()

	if err := prosvr.AddUsersToProject("validated-project-id", "validated-user-id"); err != nil {
		t.Fatal(err)
	}
}

func testProjectNotFound(t *testing.T) {
	prosvr := NewProjectController()

	if err := prosvr.AddUsersToProject("invalidated-project-id", "validated-user-id"); err.Error() != "project invalidated-project-id not found" {
		t.Fatal("want project not found")
	}
}

func TestListProjectUsers(t *testing.T) {
	t.Run("OK", testListProjectUsersOK)
	t.Run("ProjectNotFound", testListProjectUsersProjectNotFound)
}

func testListProjectUsersOK(t *testing.T) {
	prosvr := NewProjectController()

	_, err := prosvr.ListProjectUsers("validated-project-id")
	if err != nil {
		t.Fatal(err)
	}
}

func testListProjectUsersProjectNotFound(t *testing.T) {
	prosvr := NewProjectController()

	_, err := prosvr.ListProjectUsers("invalidated-project-id")
	if err.Error() != "project invalidated-project-id not found" {
		t.Fatal("want project not found")
	}
}

func TestRemoveUsersFromProject(t *testing.T) {
	t.Run("OK", testRemoveOK)
	t.Run("ProjectNotFound", testRemoveProjectNotFound)
}

func testRemoveOK(t *testing.T) {
	prosvr := NewProjectController()

	if err := prosvr.RemoveUsersFromProject("validated-project-id", "validated-user-id"); err != nil {
		t.Fatal(err)
	}
}

func testRemoveProjectNotFound(t *testing.T) {
	prosvr := NewProjectController()

	if err := prosvr.RemoveUsersFromProject("invalidated-project-id", "validated-user-id"); err.Error() != "project invalidated-project-id not found" {
		t.Fatal("want project not found")
	}
}

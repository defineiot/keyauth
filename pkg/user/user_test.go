package user_test

import (
	"testing"
)

func TestCreateUser(t *testing.T) {
	t.Run("OK", testCreateOK)
	t.Run("DomainNotFound", testCreateDomainNotFound)
}

func testCreateOK(t *testing.T) {
	usersvr := NewUserController()

	_, err := usersvr.CreateUser("validated-domain-id", "validated-username", "validated-password", "")
	if err != nil {
		t.Fatal(err)
	}
}

func testCreateDomainNotFound(t *testing.T) {
	usersvr := NewUserController()

	_, err := usersvr.CreateUser("invalidated-domain-id", "validated-username", "validated-password", "")
	if err.Error() != "domain invalidated-domain-id not exist" {
		t.Fatal("want domain not found")
	}
}

func TestGetUser(t *testing.T) {
	t.Run("OK", testGetUserOK)
	t.Run("NotFound", testGetUserNotFound)
}

func testGetUserOK(t *testing.T) {
	usersvr := NewUserController()

	_, err := usersvr.GetUser("validated-user-id")
	if err != nil {
		t.Fatal(err)
	}
}

func testGetUserNotFound(t *testing.T) {
	usersvr := NewUserController()

	_, err := usersvr.GetUser("invalidated-user-id")
	if err.Error() != "user invalidated-user-id not found" {
		t.Fatal("want not found error")
	}
}

func TestListUser(t *testing.T) {
	t.Run("OK", testListUserOK)
	t.Run("DomainNotFound", testListUserDomainNotFound)
}

func testListUserOK(t *testing.T) {
	usersvr := NewUserController()

	_, err := usersvr.ListUser("validated-domain-id")
	if err != nil {
		t.Fatal(err)
	}
}

func testListUserDomainNotFound(t *testing.T) {
	usersvr := NewUserController()

	_, err := usersvr.ListUser("invalidated-domain-id")
	if err.Error() != "domain invalidated-domain-id not exist" {
		t.Fatal("want domain not found")
	}
}

func TestSetUserDefaultProject(t *testing.T) {
	t.Run("OK", testSetDPOK)
	t.Run("ProjectNotFound", testSetProjectNotFound)
}

func testSetDPOK(t *testing.T) {
	usersvr := NewUserController()

	if err := usersvr.SetUserDefaultProject("validated-user-id", "validated-project-id"); err != nil {
		t.Fatal(err)
	}
}

func testSetProjectNotFound(t *testing.T) {
	usersvr := NewUserController()

	if err := usersvr.SetUserDefaultProject("validated-user-id", "invalidated-project-id"); err.Error() != "project invalidated-project-id not exist" {
		t.Fatal("want project not found")
	}
}

func TestListUserProjects(t *testing.T) {
	t.Run("OK", testListUPOK)
}

func testListUPOK(t *testing.T) {
	usersvr := NewUserController()

	_, err := usersvr.ListUserProjects("validated-user-id")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAddProjectsToUser(t *testing.T) {
	usersvr := NewUserController()

	if err := usersvr.AddProjectsToUser("validated-user-id", "validated-project-id"); err != nil {
		t.Fatal(err)
	}
}

func TestRemoveProjectsFromUser(t *testing.T) {
	usersvr := NewUserController()

	if err := usersvr.RemoveProjectsFromUser("validated-user-id", "validated-project-id"); err != nil {
		t.Fatal(err)
	}
}

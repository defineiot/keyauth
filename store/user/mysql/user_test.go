package mysql_test

import (
	"testing"
)

func TestUser(t *testing.T) {
	s := newTestStore()
	defer s.Close()

	user, err := s.CreateUser("test-domain-01", "test01", "123456", true, 3600, 3600)
	if err != nil {
		t.Fatal(err)
	}
	user2, err := s.CreateUser("test-domain-01", "test02", "123456", true, 3600, 3600)
	if err != nil {
		t.Fatal(err)
	}

	userByID, err := s.GetUserByID(user.ID)
	if err != nil {
		t.Fatal(err)
	}
	userByName, err := s.GetUserByName(user.DomainID, user.Name)
	if err != nil {
		t.Fatal(err)
	}
	if userByID.ID != userByName.ID {
		t.Fatal("by id user not equal by name user")
	}

	_, err = s.ValidateUser(user.DomainID, user.Name, "123456")
	if err != nil {
		t.Fatal(err)
	}

	users, err := s.ListUser(user.DomainID)
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 2 {
		t.Fatal("domain user not equal 2")
	}

	if err := s.DeleteUser(user.ID); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteUser(user2.ID); err != nil {
		t.Fatal(err)
	}
}

func TestUserProject(t *testing.T) {
	s := newTestStore()
	defer s.Close()

	user, err := s.CreateUser("test-domain-01", "test03", "123456", true, 3600, 3600)
	if err != nil {
		t.Fatal(err)
	}

	ok, err := s.CheckUserIsExistByID(user.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("create user not exist")
	}

	if err := s.AddProjectsToUser(user.ID, "project-01", "project-02"); err != nil {
		t.Fatal(err)
	}
	projects, err := s.ListUserProjects(user.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(projects) != 2 {
		t.Fatal("add 2 project but not 2")
	}

	if err := s.SetDefaultProject(user.ID, "project-01"); err != nil {
		t.Fatal(err)
	}

	if err := s.RemoveProjectsFromUser(user.ID, "project-01", "project-02"); err != nil {
		t.Fatal(err)
	}

	if err := s.DeleteUser(user.ID); err != nil {
		t.Fatal(err)
	}

}

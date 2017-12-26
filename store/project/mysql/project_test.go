package mysql_test

import (
	"testing"
)

func TestCreateProject(t *testing.T) {
	s := newTestStore()
	defer s.Close()

	p1, err := s.CreateProject("unit-test-domain-id", "test-project-01", "", true)
	if err != nil {
		t.Fatal(err)
	}
	p2, err := s.CreateProject("unit-test-domain-id", "test-project-02", "", true)
	if err != nil {
		t.Fatal(err)
	}

	if p1.Name != "test-project-01" {
		t.Fatal("project name not equal")
	}

	pGet, err := s.GetProject(p1.ID)
	if err != nil {
		t.Fatal(err)
	}
	if pGet.ID != p1.ID {
		t.Fatal("get project not equal created")
	}

	projects, err := s.ListDomainProjects("unit-test-domain-id")
	if err != nil {
		t.Fatal(err)
	}
	if len(projects) != 2 {
		t.Fatal("create 2 project in this domain, but not get 2")
	}

	if err := s.DeleteProject(p1.ID); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteProject(p2.ID); err != nil {
		t.Fatal(err)
	}
}

func TestProjectUser(t *testing.T) {
	s := newTestStore()
	defer s.Close()

	p, err := s.CreateProject("unit-test-domain-id", "test-project-01", "", true)
	if err != nil {
		t.Fatal(err)
	}

	ok, err := s.CheckProjectIsExistByID(p.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("the project not create success")
	}

	if err := s.AddUsersToProject(p.ID, "test-user-01", "test-user-02"); err != nil {
		t.Fatal(err)
	}
	uids, err := s.ListProjectUsers(p.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(uids) != 2 {
		t.Fatal("add 2 user, but not equal 2")
	}

	if err := s.RemoveUsersFromProject(p.ID, "test-user-01", "test-user-02"); err != nil {
		t.Fatal(err)
	}
	uids, err = s.ListProjectUsers(p.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(uids) != 0 {
		t.Fatal("remove 2 user, but user not equal 0")
	}

	if err := s.DeleteProject(p.ID); err != nil {
		t.Fatal(err)
	}

}

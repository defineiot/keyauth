package mysql_test

import (
	"testing"
)

func TestCreateDomain(t *testing.T) {
	s := newTestStore()
	defer s.Close()

	dom, err := s.CreateDomain("test_create_domain", "test", "", true)
	if err != nil {
		t.Fatal(err)
	}

	if dom.Name != "test_create_domain" {
		t.Fatal("domain name errror")
	}

	if err := s.DeleteDomain(dom.ID); err != nil {
		t.Fatal(err)
	}
}

func TestGetDomain(t *testing.T) {
	s := newTestStore()
	defer s.Close()

	dom, err := s.CreateDomain("test_get_domain", "test", "", true)
	if err != nil {
		t.Fatal(err)
	}

	dom1, err := s.GetDomain(dom.ID)
	if err != nil {
		t.Fatal(err)
	}

	dom2, err := s.GetDomainByName(dom.Name)
	if err != nil {
		t.Fatal(err)
	}

	if dom1.ID != dom2.ID {
		t.Fatal("get domain by id and get domain by name no equal")
	}

	if err := s.DeleteDomain(dom.ID); err != nil {
		t.Fatal(err)
	}

	ok, err := s.CheckDomainIsExistByID(dom.ID)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("the test_get_domain is not deleted")
	}
}

func TestListDomain(t *testing.T) {
	s := newTestStore()
	defer s.Close()

	dom1, err := s.CreateDomain("test_list_domain1", "test", "", true)
	if err != nil {
		t.Fatal(err)
	}

	dom2, err := s.CreateDomain("test_list_domain2", "test", "", true)
	if err != nil {
		t.Fatal(err)
	}

	domains, totalP, err := s.ListDomain(1, 1)
	if err != nil {
		if err := s.DeleteDomain(dom1.ID); err != nil {
			t.Fatal(err)
		}
		if err := s.DeleteDomain(dom2.ID); err != nil {
			t.Fatal(err)
		}
		t.Fatal(err)
	}

	_, _, err = s.ListDomain(0, 0)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(domains)
	t.Log(totalP)
	if len(domains) != 1 {
		t.Fatal("page size not equel")
	}

	if err := s.DeleteDomain(dom1.ID); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteDomain(dom2.ID); err != nil {
		t.Fatal(err)
	}

}

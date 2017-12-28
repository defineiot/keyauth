package domain_test

import (
	"testing"
)

func TestCreateDomain(t *testing.T) {
	t.Run("OK", testCreateOK)
}

func testCreateOK(t *testing.T) {
	domainsvr := NewDomainController()

	domain, err := domainsvr.CreateDomain("domain01", "unit-test", "number01", true)
	if err != nil {
		t.Fatal(err)
	}
	if domain.Name != "domain01" {
		t.Fatal("create name not correct")
	}
}

func TestListDomain(t *testing.T) {
	t.Run("OK", testListDomainOK)
}

func testListDomainOK(t *testing.T) {
	domainsvr := NewDomainController()

	domains, _, err := domainsvr.ListDomain(1, 1)
	if err != nil {
		t.Fatal(err)
	}

	if len(domains) != 1 {
		t.Fatal("need one")
	}
}

func TestGetDomin(t *testing.T) {
	t.Run("OK", testGetDomainOK)
}

func testGetDomainOK(t *testing.T) {
	domainsvr := NewDomainController()

	_, err := domainsvr.GetDomain("unit-test-id")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDestroyDomain(t *testing.T) {
	t.Run("OK", testDeleteOK)
}

func testDeleteOK(t *testing.T) {
	domainsvr := NewDomainController()

	if err := domainsvr.DestroyDomain("unit-test-id"); err != nil {
		t.Fatal(err)
	}
}

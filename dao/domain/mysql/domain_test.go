package mysql_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDomainSuit(t *testing.T) {
	suit := new(domainSuit)
	suit.SetUp()
	defer suit.TearDown()

	t.Run("CreateDomainOK", testCreateDomain(suit))
	t.Run("GetDomainByIDOK", testGetDomainByID(suit))
	t.Run("GetDomainByNameOK", testGetDomainByName(suit))
	t.Run("ListDomainOK", testListDomain(suit))
	t.Run("DeleteDomainByIDOK", testDeleteDomainByID(suit))
	t.Run("DeleteDomainByNameOK", testDeleteDomainByName(suit))
}

func testCreateDomain(s *domainSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		err := s.store.CreateDomain(s.d1)
		should.NoError(err)
		err = s.store.CreateDomain(s.d2)
		should.NoError(err)

		t.Logf("create domain(%s) success: %s", s.d1.Name, s.d1)
		t.Logf("create domain(%s) success: %s", s.d2.Name, s.d2)
	}
}

func testGetDomainByID(s *domainSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		dom1, err := s.store.GetDomainByID(s.d1.ID)
		should.NoError(err)

		t.Logf("get domain by id success: %s", dom1)
	}
}

func testGetDomainByName(s *domainSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		dom2, err := s.store.GetDomainByName(s.d2.Name)
		should.NoError(err)

		t.Logf("get domain by id success: %s", dom2)
	}
}

func testListDomain(s *domainSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		doms, total, err := s.store.ListDomain(1, 2)
		should.NoError(err)

		t.Logf("list domain success: domains: %s, total_page: %d", doms, total)
	}
}

func testDeleteDomainByID(s *domainSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		err := s.store.DeleteDomainByID(s.d1.ID)
		should.NoError(err)

		t.Logf("delete domain by id success (%s)", s.d1.ID)
	}
}

func testDeleteDomainByName(s *domainSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		err := s.store.DeleteDomainByName(s.d2.Name)
		should.NoError(err)

		t.Logf("delete domain by name success (%s)", s.d2.Name)
	}
}

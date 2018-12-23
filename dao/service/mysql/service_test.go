package mysql_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProjectSuit(t *testing.T) {
	suit := new(serviceSuit)
	suit.SetUp()
	defer suit.TearDown()

	t.Run("CreateServiceOK", testCreateServiceOK(suit))
	t.Run("GetServiceByIDOK", testGetServiceByIDOK(suit))
	t.Run("GetServiceByClientOK", testGetServiceByClientOK(suit))
	t.Run("ListServiceOK", testListServiceOK(suit))
	t.Run("RegisteServiceFeaturesOK", testRegisteServiceFeaturesOK(suit))
	t.Run("AssociateFeaturesToRole", testAssociateFeaturesToRole(suit))
	t.Run("ListServiceFeaturesOK", testListServiceFeaturesOK(suit))
	t.Run("ListRoleFeatruresOK", testListRoleFeatruresOK(suit))
	t.Run("UnlinkFeatureFromRole", testUnlinkFeatureFromRoleOK(suit))
	t.Run("DeleteServiceOK", testDeleteServiceOK(suit))
}

func testCreateServiceOK(s *serviceSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		err := s.store.CreateService(s.svr)
		should.NoError(err)

		t.Logf("create service(%s) success: %s", s.svr.Name, s.svr)
	}
}

func testGetServiceByIDOK(s *serviceSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		svr, err := s.store.GetServiceByID(s.svr.ID)
		should.NoError(err)

		t.Logf("get service by id(%s) success: %s", s.svr.Name, svr)
	}
}

func testGetServiceByClientOK(s *serviceSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		svr, err := s.store.GetServiceByClientID(s.svr.ClientID)
		should.NoError(err)

		t.Logf("get service by client(%s) success: %s", s.svr.Name, svr)
	}
}

func testListServiceOK(s *serviceSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		svrs, err := s.store.ListServices()
		should.NoError(err)

		t.Logf("list services success: %s", svrs)
	}
}

func testRegisteServiceFeaturesOK(s *serviceSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		err := s.store.RegistryServiceFeatures(s.svr.ID, s.features...)
		should.NoError(err)

		t.Logf("registe service(%s) features: %v", s.svr.Name, s.features)
	}
}

func testAssociateFeaturesToRole(s *serviceSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		err := s.store.AssociateFeaturesToRole(s.roleID)
		should.NoError(err)

		t.Logf("delete service(%s) success: %s", s.svr.Name, s.svr)
	}
}

func testListServiceFeaturesOK(s *serviceSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		features, err := s.store.ListServiceFeatures(s.svr.ID)
		should.NoError(err)

		t.Logf("list service(%s) features: %v", s.svr.Name, features)
	}
}

func testListRoleFeatruresOK(s *serviceSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		features, err := s.store.ListRoleFeatures(s.roleID)
		should.NoError(err)

		t.Logf("list role(%s) features: %v", s.roleID, features)
	}
}

func testUnlinkFeatureFromRoleOK(s *serviceSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		err := s.store.UnlinkFeatureFromRole(s.roleID, s.features...)
		should.NoError(err)

		t.Logf("unlink role(%s) features: %v", s.roleID, s.features)
	}
}

func testDeleteServiceOK(s *serviceSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		err := s.store.DeleteService(s.svr.ID)
		should.NoError(err)

		t.Logf("delete service(%s) success: %s", s.svr.Name, s.svr)
	}
}

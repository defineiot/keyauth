package etcd_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInstanceSuit(t *testing.T) {
	suit := new(instanceSuit)
	suit.SetUp()
	defer suit.TearDown()

	t.Run("CreateInstanceOK", testCreateInstanceOK(suit))
	t.Run("ListServiceInstancesOK", testListServiceInstanceOK(suit))
	t.Run("testDeleteServiceInstanceOK", testDeleteServiceInstanceOK(suit))
}

func testCreateInstanceOK(s *instanceSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		err := s.store.RegistryServiceInstance(s.ins)
		should.NoError(err)

		t.Logf("create instance(%s) success: %s", s.ins.Name, s.ins)
	}
}

func testListServiceInstanceOK(s *instanceSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		inss, err := s.store.ListServiceInstances(s.ins.ServiceID)
		should.NoError(err)

		t.Logf("list service(%s) instance success: %s", s.ins.ServiceID, inss)
	}
}

func testDeleteServiceInstanceOK(s *instanceSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		err := s.store.UnRegisteServiceInstance(s.ins.ServiceID, s.ins.ID)
		should.NoError(err)

		t.Logf("delete service(%s) instance success: %s", s.ins.ServiceID, s.ins)
	}
}

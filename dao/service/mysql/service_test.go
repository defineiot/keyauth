package mysql_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/defineiot/keyauth/dao/service"
)

var serviceName string

func TestService(t *testing.T) {
	t.Run("CreateServiceOK", testCreateServiceOK)
	t.Run("ListServiceOK", testListServiceOK)
	t.Run("GetServiceOK", testListServiceOK)
	t.Run("RegistryFeatures", testRegistryServiceFeatures)
	t.Run("ListServiceFeatures", testListServiceFeatures)
	t.Run("DeleteServiceOK", testDeleteServiceOK)
}

func testCreateServiceOK(t *testing.T) {
	s := newTestStore()
	defer s.Close()

	svr, err := s.CreateService("test_service02", "just for unit test", "ooxxooxxooxx")
	assert.NoError(t, err)
	assert.NotNil(t, svr)
	assert.Equal(t, "test_service02", svr.Name)

	serviceName = svr.Name
}

func testListServiceOK(t *testing.T) {
	s := newTestStore()
	defer s.Close()

	services, err := s.ListServices()
	assert.NoError(t, err)
	assert.NotEmpty(t, len(services))
}

func testGetServiceOK(t *testing.T) {
	s := newTestStore()
	defer s.Close()

	svr, err := s.GetService(serviceName)
	assert.NoError(t, err)
	assert.NotNil(t, svr)
	assert.Equal(t, "test_service02", svr.Name)
}

func testRegistryServiceFeatures(t *testing.T) {
	s := newTestStore()
	defer s.Close()

	f1 := service.Feature{Name: "F1", Method: "GET", Endpoint: "/v1/F1"}
	f2 := service.Feature{Name: "F2", Method: "POST", Endpoint: "/v1/F2"}
	f3 := service.Feature{Name: "F3", Method: "DELETE", Endpoint: "/v1/F3"}
	err := s.RegistryServiceFeatures(serviceName, f1, f2, f3)
	assert.NoError(t, err)
}

func testListServiceFeatures(t *testing.T) {
	s := newTestStore()
	defer s.Close()

	features, err := s.ListServiceFeatures(serviceName)
	assert.NoError(t, err)
	assert.NotEmpty(t, len(features))
}

func testDeleteServiceOK(t *testing.T) {
	s := newTestStore()
	defer s.Close()

	err := s.DeleteService(serviceName)
	assert.NoError(t, err)
}

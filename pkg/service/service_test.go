package service_test

import (
	"testing"
)

func TestCreateService(t *testing.T) {
	t.Run("OK", testCreateOK)
}

func testCreateOK(t *testing.T) {
	svr := NewServiceController()

	if _, err := svr.CreateService("validated-name", "unit-test"); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteService(t *testing.T) {
	t.Run("OK", testDeleteOK)
}

func testDeleteOK(t *testing.T) {
	svr := NewServiceController()

	if err := svr.DeleteService("validated-sid"); err != nil {
		t.Fatal(err)
	}
}

func TestListService(t *testing.T) {
	t.Run("OK", testListOK)
}

func testListOK(t *testing.T) {
	svr := NewServiceController()

	if _, err := svr.ListService(); err != nil {
		t.Fatal(err)
	}
}

func TestGetService(t *testing.T) {
	t.Run("OK", testGetOK)
}

func testGetOK(t *testing.T) {
	svr := NewServiceController()

	if _, err := svr.GetService("validated-sid"); err != nil {
		t.Fatal(err)
	}

}

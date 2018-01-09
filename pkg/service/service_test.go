package service_test

import (
	"testing"
)

func TestCreateService(t *testing.T) {
	svr := NewServiceController()

	if _, err := svr.CreateService("validated-name", "unit-test"); err != nil {
		t.Fatal(err)
	}
}

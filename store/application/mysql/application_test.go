package mysql_test

import (
	"testing"
)

func TestRegistration(t *testing.T) {
	s := newTestStore()
	defer s.Close()

	app1, err := s.Registration("unit-test-user-id", "unit-test-app1", "http://127.0.0.1/test", "public", "", "")
	if err != nil {
		t.Fatal(err)
	}

	app2, err := s.Registration("unit-test-user-id", "unit-test-app2", "http://127.0.0.1/test", "public", "", "")
	if err != nil {
		t.Fatal(err)
	}

	if app1.Name != "unit-test-app1" {
		t.Fatal("app name not equal")
	}
	if app2.Name != "unit-test-app2" {
		t.Fatal("app name not equal")
	}

	if app1.Client == nil {
		t.Fatal("app's client not be created")
	}
	client, err := s.GetClient(app1.Client.ClientID)
	if err != nil {
		t.Fatal(err)
	}
	if client.ClientID == "" {
		t.Fatal("client id is \"\"")
	}

	apps, err := s.GetUserApps("unit-test-user-id")
	if err != nil {
		t.Fatal(err)
	}

	if len(apps) != 2 {
		t.Fatal("the user app not equal 2")
	}

	if err := s.Unregistration(app1.ID); err != nil {
		t.Fatal(err)
	}
	if err := s.Unregistration(app2.ID); err != nil {
		t.Fatal(err)
	}
}

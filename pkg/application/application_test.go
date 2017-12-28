package application_test

import (
	"testing"
)

func TestRegisteAPP(t *testing.T) {
	t.Run("OK", testRegisteSuccess)
	t.Run("UserNotFound", testRegisteUserNotFound)
}

func testRegisteSuccess(t *testing.T) {

	appsvr := NewAppController()

	_, err := appsvr.RegisterApplication("unit-exist-user", "app01", "http://baidu.com", "public", "", "")
	if err != nil {
		t.Fatal(err)
	}
}

func testRegisteUserNotFound(t *testing.T) {
	appsvr := NewAppController()
	_, err := appsvr.RegisterApplication("unit-not-exist-user", "app01", "http://baidu.com", "public", "", "")
	if err == nil {
		t.Fatal("want not find user error")
	}
	if err.Error() != "user unit-not-exist-user not exist" {
		t.Fatal("test user unit-not-exist-user not find ,but find user")
	}
}

func TestUnRegisteAPP(t *testing.T) {
	t.Run("OK", testUnRegisteAPPOK)
	t.Run("NotFound", testUnRegisteAPPNotFound)
}

func testUnRegisteAPPOK(t *testing.T) {
	appsvr := NewAppController()
	if err := appsvr.UnregisteApplication("unit-test-app-id"); err != nil {
		t.Fatal(err)
	}
}

func testUnRegisteAPPNotFound(t *testing.T) {
	appsvr := NewAppController()
	if err := appsvr.UnregisteApplication("unit-test-app-id-not-find"); err == nil {
		t.Fatal("want not find error")
	}
}

func TestGetUserAPP(t *testing.T) {
	t.Run("OK", testGetUserAPPOK)
	t.Run("NotFound", testGetUserAPPNotFound)
}

func testGetUserAPPOK(t *testing.T) {
	appsvr := NewAppController()
	if apps, err := appsvr.GetUserApplications("unit-exist-user"); err != nil {
		t.Fatal(err)
	} else {
		if len(apps) != 1 {
			t.Fatal("want one app")
		}
	}
}

func testGetUserAPPNotFound(t *testing.T) {
	appsvr := NewAppController()
	if _, err := appsvr.GetUserApplications("unit-test-app-id-not-find"); err == nil {
		t.Fatal("want user not found error")
	}
}

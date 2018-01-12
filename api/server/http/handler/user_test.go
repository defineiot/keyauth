package handler_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"openauth/api/server/http/handler"
	"strings"
	"testing"

	"github.com/json-iterator/go"
)

var (
	userID string
	appID  string
)

func TestCreateUser(t *testing.T) {
	t.Run("OK", testCreateUserOK)
}

func testCreateUserOK(t *testing.T) {
	payload := strings.NewReader(`{"name": "unit-test", "password": "unit-test"}`)

	req, err := http.NewRequest("POST", "/v1/users/", payload)
	if err != nil {
		t.Fatal(err)
	}

	rr := newRequestRecorder(req, "POST", "/v1/users/", handler.CreateUser)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(rr.Result().Body)
	if err != nil {
		t.Errorf("read body data error, %s", err)
	}

	userID = jsoniter.Get(body, "data", "id").ToString()
}

func TestRegisteApplication(t *testing.T) {
	t.Run("OK", testRegisteOK)
}

func testRegisteOK(t *testing.T) {
	if userID == "" {
		t.Fatal("create not save user id")
	}

	payload := strings.NewReader(`{"name": "unit-test01", "client_type": "public"}`)

	url := fmt.Sprintf("/v1/users/%s/applications/", userID)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		t.Fatal(err)
	}

	rr := newRequestRecorder(req, "POST", "/v1/users/:uid/applications/", handler.RegisteApplication)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(rr.Result().Body)
	if err != nil {
		t.Errorf("read body data error, %s", err)
	}

	appID = jsoniter.Get(body, "data", "id").ToString()
}

func TestGetUserApplications(t *testing.T) {
	t.Run("OK", testGetUserAPPOK)
}

func testGetUserAPPOK(t *testing.T) {
	if appID == "" {
		t.Fatal("create not save app id")
	}
	url := fmt.Sprintf("/v1/users/%s/applications/%s/", userID, appID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := newRequestRecorder(req, "GET", "/v1/users/:uid/applications/:aid/", handler.GetUserApplications)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestUnRegisteApplication(t *testing.T) {
	t.Run("OK", testUnRegisteOK)
}

func testUnRegisteOK(t *testing.T) {
	if appID == "" {
		t.Fatal("create not save app id")
	}
	url := fmt.Sprintf("/v1/users/%s/applications/%s/", userID, appID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := newRequestRecorder(req, "DELETE", "/v1/users/:uid/applications/:aid/", handler.UnRegisteApplication)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}

func TestDeleteUser(t *testing.T) {
	t.Run("OK", testDeleteUserOK)
}

func testDeleteUserOK(t *testing.T) {
	if userID == "" {
		t.Fatal("create not save user id")
	}
	url := fmt.Sprintf("/v1/users/%s/", userID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := newRequestRecorder(req, "DELETE", "/v1/users/:uid/", handler.DeleteUser)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}

package handler_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/json-iterator/go"

	"openauth/api/server/http/handler"
)

var (
	serviceID string
)

func TestCreateService(t *testing.T) {
	t.Run("CreateOK", testCreateOK)

}

func testCreateOK(t *testing.T) {
	payload := strings.NewReader(`{"name": "test"}`)

	req, err := http.NewRequest("POST", "/v1/services/", payload)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	handler.CreateService(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(rr.Result().Body)
	if err != nil {
		t.Errorf("read body data error, %s", err)
	}

	serviceID = jsoniter.Get(body, "data", "id").ToString()
}

func TestListService(t *testing.T) {
	t.Run("OK", testListOK)
}

func testListOK(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/services/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := newRequestRecorder(req, "GET", "/v1/services/", handler.ListService)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestGetService(t *testing.T) {
	t.Run("OK", testGetOK)
}

func testGetOK(t *testing.T) {
	if serviceID == "" {
		t.Fatal("create not save service id")
	}
	url := fmt.Sprintf("/v1/services/%s/", serviceID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := newRequestRecorder(req, "GET", "/v1/services/:sid/", handler.GetService)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestDeleteService(t *testing.T) {
	t.Run("DeleteOK", testDeleteOK)
}

func testDeleteOK(t *testing.T) {
	if serviceID == "" {
		t.Fatal("create not save service id")
	}
	url := fmt.Sprintf("/v1/services/%s/", serviceID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := newRequestRecorder(req, "DELETE", "/v1/services/:sid/", handler.DeleteService)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}

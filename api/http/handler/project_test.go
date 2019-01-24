package handler_test

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// 	"strings"
// 	"testing"

// 	"github.com/json-iterator/go"
// 	"github.com/stretchr/testify/assert"

// 	"github.com/defineiot/keyauth/api/global"
// 	"github.com/defineiot/keyauth/api/http/handler"
// )

// var (
// 	projectID       string
// 	projectDomainID string
// )

// func TestProject(t *testing.T) {
// 	t.Run("CreateOK", testCreateProjectOK)
// 	t.Run("GetOK", testGetProjectOK)
// 	t.Run("ListOK", testListProjectOK)
// 	t.Run("DeleteOK", testDeleteProjectOK)
// }

// func testCreateProjectOK(t *testing.T) {
// 	dom, err := global.Store.CreateDomain("project-domain01", "", "test", true)
// 	assert.NoError(t, err)
// 	projectDomainID = dom.ID

// 	payload := strings.NewReader(`{"name": "unit-test-project01", "display_name": "simple project"}`)

// 	global.Store.DeleteProjectByName("unit-test-project01", projectDomainID)

// 	url := fmt.Sprintf("/v1/domains/%s/projects/", projectDomainID)
// 	req, err := http.NewRequest("POST", url, payload)
// 	assert.NoError(t, err)

// 	rr := newRequestRecorder(req, "POST", "/v1/domains/:did/projects/", handler.CreateProject)

// 	if status := rr.Code; status != http.StatusCreated {
// 		msg, _ := ioutil.ReadAll(rr.Result().Body)
// 		t.Errorf("handler returned wrong status code: got %v want %v, msg: %s", status, http.StatusCreated, string(msg))
// 	}

// 	body, err := ioutil.ReadAll(rr.Result().Body)
// 	assert.NoError(t, err, "read body data error, %s", err)

// 	projectName := jsoniter.Get(body, "data", "name").ToString()
// 	assert.Equal(t, "unit-test-project01", projectName)

// 	projectID = jsoniter.Get(body, "data", "id").ToString()
// }

// func testGetProjectOK(t *testing.T) {
// 	url := fmt.Sprintf("/v1/projects/%s/", projectID)
// 	req, err := http.NewRequest("GET", url, nil)
// 	assert.NoError(t, err)

// 	rr := newRequestRecorder(req, "GET", "/v1/projects/:pid/", handler.GetProject)

// 	if status := rr.Code; status != http.StatusOK {
// 		msg, _ := ioutil.ReadAll(rr.Result().Body)
// 		t.Errorf("handler returned wrong status code: got %v want %v, msg: %s", status, http.StatusOK, string(msg))
// 	}

// 	body, err := ioutil.ReadAll(rr.Result().Body)
// 	assert.NoError(t, err, "read body data error, %s", err)

// 	projectName := jsoniter.Get(body, "data", "name").ToString()
// 	assert.Equal(t, "unit-test-project01", projectName)

// }

// func testListProjectOK(t *testing.T) {
// 	url := fmt.Sprintf("/v1/domains/%s/projects/", projectDomainID)
// 	req, err := http.NewRequest("GET", url, nil)
// 	assert.NoError(t, err)

// 	rr := newRequestRecorder(req, "GET", "/v1/domains/:did/projects/", handler.ListProject)

// 	if status := rr.Code; status != http.StatusOK {
// 		msg, _ := ioutil.ReadAll(rr.Result().Body)
// 		t.Errorf("handler returned wrong status code: got %v want %v, msg: %s", status, http.StatusOK, string(msg))
// 	}

// 	body, err := ioutil.ReadAll(rr.Result().Body)
// 	assert.NoError(t, err, "read body data error, %s", err)

// 	projectName := jsoniter.Get(body, "data", 0, "name").ToString()
// 	assert.Equal(t, "unit-test-project01", projectName)
// }

// func testDeleteProjectOK(t *testing.T) {
// 	url := fmt.Sprintf("/v1/projects/%s/", projectID)
// 	req, err := http.NewRequest("DELETE", url, nil)
// 	assert.NoError(t, err)

// 	rr := newRequestRecorder(req, "DELETE", "/v1/projects/:pid/", handler.DeleteProject)

// 	if status := rr.Code; status != http.StatusNoContent {
// 		msg, _ := ioutil.ReadAll(rr.Result().Body)
// 		t.Errorf("handler returned wrong status code: got %v want %v, msg: %s", status, http.StatusNoContent, string(msg))
// 	}

// 	err = global.Store.DeleteDomain("id", projectDomainID)
// 	assert.NoError(t, err)
// }

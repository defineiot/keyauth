package handler_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"

	"github.com/defineiot/keyauth/api/global"
	"github.com/defineiot/keyauth/api/http/handler"
	"github.com/defineiot/keyauth/api/http/router"
	"github.com/defineiot/keyauth/internal/cache"
	"github.com/defineiot/keyauth/internal/conf/mock"
	"github.com/defineiot/keyauth/store"
	domain "github.com/defineiot/keyauth/store/domain/mysql"
	project "github.com/defineiot/keyauth/store/project/mysql"
)

var (
	domainID string
)

// Mocks a handler and returns a httptest.ResponseRecorder
func newRequestRecorder(req *http.Request, method string, path string, handleFunc http.HandlerFunc) *httptest.ResponseRecorder {
	r := router.NewRouter()
	r.HandlerFunc(method, path, handleFunc)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	r.Router.ServeHTTP(rr, req)
	return rr
}

func TestDomain(t *testing.T) {
	t.Run("CreateOK", testCreateDomainOK)
	t.Run("GetOK", testGetDomainOK)
	t.Run("ListOK", testListDomainOK)
	t.Run("DeleteOK", testDeleteDomainOK)
}

func testCreateDomainOK(t *testing.T) {

	payload := strings.NewReader(`{"name": "unit-test-domain01", "display_name": "test"}`)

	global.Store.DeleteDomain("name", "unit-test-domain01")

	req, err := http.NewRequest("POST", "/v1/domains/", payload)
	assert.NoError(t, err)

	rr := newRequestRecorder(req, "POST", "/v1/domains/", handler.CreateDomain)

	if status := rr.Code; status != http.StatusCreated {
		msg, _ := ioutil.ReadAll(rr.Result().Body)
		t.Errorf("handler returned wrong status code: got %v want %v, msg: %s", status, http.StatusCreated, string(msg))
	}

	body, err := ioutil.ReadAll(rr.Result().Body)
	assert.NoError(t, err, "read body data error, %s", err)

	domainName := jsoniter.Get(body, "data", "name").ToString()
	assert.Equal(t, "unit-test-domain01", domainName)

	domainID = jsoniter.Get(body, "data", "id").ToString()
}

func testGetDomainOK(t *testing.T) {
	uri := fmt.Sprintf("/v1/domains/%s/", domainID)

	req, err := http.NewRequest("GET", uri, nil)
	assert.NoError(t, err)

	rr := newRequestRecorder(req, "GET", "/v1/domains/:did/", handler.GetDomain)

	if status := rr.Code; status != http.StatusOK {
		msg, _ := ioutil.ReadAll(rr.Result().Body)
		t.Errorf("handler returned wrong status code: got %v want %v, msg: %s", status, http.StatusOK, string(msg))
	}

	body, err := ioutil.ReadAll(rr.Result().Body)
	assert.NoError(t, err, "read body data error, %s", err)

	domainName := jsoniter.Get(body, "data", "name").ToString()
	assert.Equal(t, "unit-test-domain01", domainName)
}

func testListDomainOK(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/domains/", nil)
	assert.NoError(t, err)

	rr := newRequestRecorder(req, "GET", "/v1/domains/", handler.ListDomain)

	if status := rr.Code; status != http.StatusOK {
		msg, _ := ioutil.ReadAll(rr.Result().Body)
		t.Errorf("handler returned wrong status code: got %v want %v, msg: %s", status, http.StatusOK, string(msg))
	}

	body, err := ioutil.ReadAll(rr.Result().Body)
	assert.NoError(t, err, "read body data error, %s", err)

	domainName := jsoniter.Get(body, "data", 0, "name").ToString()
	assert.Equal(t, "unit-test-domain01", domainName)
}

func testDeleteDomainOK(t *testing.T) {
	uri := fmt.Sprintf("/v1/domains/%s/", domainID)

	req, err := http.NewRequest("DELETE", uri, nil)
	assert.NoError(t, err)

	rr := newRequestRecorder(req, "DELETE", "/v1/domains/:did/", handler.DeleteDomain)

	if status := rr.Code; status != http.StatusNoContent {
		msg, _ := ioutil.ReadAll(rr.Result().Body)
		t.Errorf("handler returned wrong status code: got %v want %v, msg: %s", status, http.StatusNoContent, string(msg))
	}
}

func init() {
	conf := mock.NewConfig()

	log, err := conf.GetLogger()
	if err != nil {
		panic(err)
	}
	db, err := conf.GetDBConn()
	if err != nil {
		panic(err)
	}
	dom, err := domain.NewDomainStore(db)
	if err != nil {
		panic(err)
	}
	pro, err := project.NewProjectStore(db)
	if err != nil {
		panic(err)
	}

	store := store.NewStore(dom, log, pro)
	store.SetCache(cache.Newmemcache(1000), time.Minute*5)

	global.Conf = conf
	global.Log = log
	global.Store = store
}

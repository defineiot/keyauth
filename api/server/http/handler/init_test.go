package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"openauth/api/config/mock"
	"openauth/api/server/http/handler"
	"openauth/api/server/http/router"
)

func TestInitController(t *testing.T) {
	conf := mock.NewConfig()

	if err := handler.InitController(conf); err != nil {
		t.Fatal(err)
	}
}

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

func init() {
	conf := mock.NewConfig()

	if err := handler.InitController(conf); err != nil {
		panic(err)
	}
}

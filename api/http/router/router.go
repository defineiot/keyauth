package router

import (
	"net/http"
	"openauth/api/http/response"
	"reflect"
	"runtime"
	"strings"

	"openauth/api/http/context"

	"github.com/julienschmidt/httprouter"
)

// MyRouter is an hack for httprouter
type MyRouter struct {
	Router      *httprouter.Router
	v1endpoints map[string]map[string]string
}

// NewRouter use to new an router
func NewRouter() *MyRouter {
	hrouter := &httprouter.Router{
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      true,
		HandleMethodNotAllowed: true,
		HandleOPTIONS:          true,
	}

	hrouter.NotFound = new(notFoundHandler)
	hrouter.MethodNotAllowed = new(methodNotAllowedHandler)

	ep := make(map[string]map[string]string)

	r := &MyRouter{Router: hrouter, v1endpoints: ep}

	return r
}

// Handler is an adapter which allows the usage of an http.Handler as a
// request handle.
func (r *MyRouter) Handler(method, path string, handler http.Handler) {
	r.Router.Handle(method, path,
		func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
			context.SetParamsToContext(req, w, handler, ps)
		},
	)
}

// HandlerFunc yes hack
func (r *MyRouter) HandlerFunc(method, path string, handleFunc http.HandlerFunc) {

	fn := runtime.FuncForPC(reflect.ValueOf(handleFunc).Pointer()).Name()
	feathure := strings.Split(fn, ".")[1]

	if strings.HasPrefix(path, "/v1/") {
		_, ok := r.v1endpoints[method]
		if !ok {
			mm := make(map[string]string)
			r.v1endpoints[method] = mm
		}
		r.v1endpoints[method][feathure] = path
	}

	r.Handler(method, path, http.HandlerFunc(handleFunc))
}

// AddV1Root add root to api
func (r *MyRouter) AddV1Root() {
	r.Router.Handle("GET", "/v1/",
		func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
			response.Success(w, http.StatusOK, r.v1endpoints)
			return
		},
	)
}

// GetEndpoints get router's fl
func (r *MyRouter) GetEndpoints() map[string]map[string]string {
	return r.v1endpoints
}

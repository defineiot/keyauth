package router

import (
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"openauth/api/http/context"

	"github.com/julienschmidt/httprouter"
)

// MyRouter is an hack for httprouter
type MyRouter struct {
	Router    *httprouter.Router
	endpoints map[string]string
}

// NewRouter use to new an router
func NewRouter() *MyRouter {
	hrouter := &httprouter.Router{
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      true,
		HandleMethodNotAllowed: true,
		HandleOPTIONS:          true,
	}
	ep := make(map[string]string)

	r := &MyRouter{Router: hrouter, endpoints: ep}

	return r
}

// Handler is an adapter which allows the usage of an http.Handler as a
// request handle.
func (r *MyRouter) Handler(method, path string, handler http.Handler, name string) {
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

	r.endpoints[path] = feathure

	r.Handler(method, path, http.HandlerFunc(handleFunc), feathure)
}

// GetEndpoints get router's fl
func (r *MyRouter) GetEndpoints() map[string]string {
	return r.endpoints
}

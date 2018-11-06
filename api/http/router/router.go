package router

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/defineiot/keyauth/api/global"
	"github.com/defineiot/keyauth/api/http/context"
	"github.com/defineiot/keyauth/api/http/response"
	"github.com/defineiot/keyauth/dao/token"
	"github.com/defineiot/keyauth/internal/exception"
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
	}

	hrouter.NotFound = notFoundHandler
	hrouter.MethodNotAllowed = methodNotAllowedHandler

	ep := make(map[string]map[string]string)

	r := &MyRouter{Router: hrouter, v1endpoints: ep}

	return r
}

// Handler is an adapter which allows the usage of an http.Handler as a
// request handle.
func (r *MyRouter) Handler(method, path, featureName string, handler http.Handler) {
	r.Router.Handle(method, path,
		func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
			rctx := new(context.ReqContext)
			switch featureName {
			case "IssueToken", "ValidateToken", "RevolkToken", "IssueVerifyCode", "RegistryUser":
			case "RegistryServiceFeatures":
				authHeader := req.Header.Get("Authorization")
				if authHeader == "" {
					response.Failed(w, exception.NewUnauthorized("Authorization missed in header"))
					return
				}

				headerSlice := strings.Split(authHeader, " ")
				if len(headerSlice) != 2 {
					response.Failed(w, exception.NewUnauthorized("Authorization header value is not validated, must be: {token_type} {token}"))
					return
				}

				access := headerSlice[1]

				t, err := global.Store.ValidateToken(access)
				if err != nil {
					response.Failed(w, err)
					return
				}

				if t.GrantType != token.CLIENT {
					response.Failed(w, exception.NewForbidden("feature: RegistryServiceFeatures only for client grant"))
					return
				}

				if t.ServiceName != ps.ByName("sn") {
					response.Failed(w, exception.NewForbidden("the token: %s is not belong this service: %s", t.AccessToken, t.ServiceName))
					return
				}

			default:
				var hasPerm bool
				authHeader := req.Header.Get("Authorization")
				if authHeader == "" {
					response.Failed(w, exception.NewUnauthorized("Authorization missed in header"))
					return
				}

				headerSlice := strings.Split(authHeader, " ")
				if len(headerSlice) != 2 {
					response.Failed(w, exception.NewUnauthorized("Authorization header value is not validated, must be: {token_type} {token}"))
					return
				}

				access := headerSlice[1]

				t, err := global.Store.ValidateToken(access)
				if err != nil {
					response.Failed(w, err)
					return
				}

				if t.GrantType == token.CLIENT {
					response.Failed(w, exception.NewForbidden("client grant token only for RegistryServiceFeatures"))
					return
				}

				u, err := global.Store.GetUser(t.DomainID, t.UserID)
				if err != nil {
					response.Failed(w, exception.NewUnauthorized(err.Error()))
					return
				}

				for _, rn := range u.RoleNames {
					switch rn {
					case "system_admin":
						t.IsSystemAdmin = true
					case "domain_admin":
						t.IsDomainAdmin = true
					default:
					}
				}

				for _, rn := range u.RoleNames {
					if rn == "system_admin" {
						hasPerm = true
						break
					}

					role, err := global.Store.GetRole(rn)
					if err != nil {
						response.Failed(w, exception.NewUnauthorized(err.Error()))
						return
					}

					for _, f := range role.Featrues {
						if f.Name == featureName {
							hasPerm = true
						}
					}
				}

				if !hasPerm {
					response.Failed(w, exception.NewForbidden("your has no permission to access feature: %s", featureName))
					return
				}

				rctx.Token = t
			}

			rctx.PS = ps
			context.SetReqContext(req, w, handler, rctx)
		},
	)

}

// HandlerFunc yes hack
func (r *MyRouter) HandlerFunc(method, path, featureName string, handleFunc http.HandlerFunc) {
	// fn := runtime.FuncForPC(reflect.ValueOf(handleFunc).Pointer()).Name()
	// feathure := strings.Split(fn, ".")[1]

	if strings.HasPrefix(path, "/v1/") {
		_, ok := r.v1endpoints[method]
		if !ok {
			mm := make(map[string]string)
			r.v1endpoints[method] = mm
		}
		if featureName != "" {
			r.v1endpoints[method][featureName] = path
		}
	}

	r.Handler(method, path, featureName, http.HandlerFunc(handleFunc))
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

package context

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/defineiot/keyauth/dao/models"
)

// ReqContextKey req context key
const ReqContextKey = "request-context-key"

// ReqContext context
type ReqContext struct {
	PS    httprouter.Params
	Token *models.Token
}

// SetReqContext use to get context
func SetReqContext(req *http.Request, w http.ResponseWriter, next http.Handler, rctx *ReqContext) {
	ctx := context.WithValue(req.Context(), ReqContextKey, rctx)
	req = req.WithContext(ctx)
	next.ServeHTTP(w, req)
}

// GetParamsFromContext use to get httprouter ps from context
func GetParamsFromContext(req *http.Request) httprouter.Params {
	return req.Context().Value(ReqContextKey).(*ReqContext).PS
}

// GetTokenFromContext use to get httprouter ps from context
func GetTokenFromContext(req *http.Request) *models.Token {
	return req.Context().Value(ReqContextKey).(*ReqContext).Token
}

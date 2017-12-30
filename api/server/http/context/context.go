package context

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const httprouterParamsKey = "params"

// SetParamsToContext use to set http ps into context
func SetParamsToContext(req *http.Request, w http.ResponseWriter, handler http.Handler, ps httprouter.Params) {
	ctx := context.WithValue(req.Context(), httprouterParamsKey, ps)
	handler.ServeHTTP(w, req.WithContext(ctx))
}

// GetParamsFromContext use to get httprouter ps from context
func GetParamsFromContext(req *http.Request) httprouter.Params {
	return req.Context().Value(httprouterParamsKey).(httprouter.Params)
}

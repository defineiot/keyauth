package response

import (
	"encoding/json"
	"net/http"

	"github.com/defineiot/keyauth/internal/exception"
)

// Response to be used by controllers.
type Response struct {
	Code     int         `json:"code"`
	Message  string      `json:"message,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	TotalPag int64       `json:"total_page,omitempty"`
}

// Failed use to response error messge
func Failed(w http.ResponseWriter, err error) {
	msg := err.Error()
	httCode := http.StatusInternalServerError
	customCode := 0

	switch t := err.(type) {
	case *exception.BadRequest:
		httCode = t.Code()
		customCode = 4000
	case *exception.NotFound:
		httCode = t.Code()
		customCode = 4004
	case *exception.InternalServerError:
		httCode = t.Code()
		customCode = 5000
	case *exception.Unauthorized:
		httCode = t.Code()
		customCode = 4001
	case *exception.MethodNotAllowed:
		httCode = t.Code()
		customCode = 4005
	case *exception.Forbidden:
		httCode = t.Code()
		customCode = 4003
	case *exception.Expired:
		httCode = t.Code()
		customCode = 4100
	default:
		httCode = http.StatusInternalServerError
	}

	resp := Response{
		Code:    customCode,
		Message: msg,
	}

	// set response heanders
	w.Header().Set("Content-Type", "application/json")

	// if marshal json error, use string to response
	respByt, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error", "message": "encoding to json error"}`))
		return
	}

	w.WriteHeader(httCode)
	w.Write(respByt)
	return
}

// Success use to response success data
func Success(w http.ResponseWriter, code int, data interface{}) {
	resp := Response{
		Message: "",
		Data:    data,
	}

	// set response heanders
	w.Header().Set("Content-Type", "application/json")

	respByt, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error", "message": "encoding to json error"}`))
		return
	}

	w.WriteHeader(code)
	w.Write(respByt)
	return
}

// SuccessWithPage use to response success data
func SuccessWithPage(w http.ResponseWriter, code int, data interface{}, totoalPage int64) {
	resp := Response{
		Message:  "",
		Data:     data,
		TotalPag: totoalPage,
	}

	// set response heanders
	w.Header().Set("Content-Type", "application/json")

	respByt, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error", "message": "encoding to json error"}`))
		return
	}

	w.WriteHeader(code)
	w.Write(respByt)
	return
}

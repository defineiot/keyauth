package response

import (
	"encoding/json"
	"net/http"
	"openauth/api/exception"
)

// Response to be used by controllers.
type Response struct {
	Status   string      `json:"status"`
	Message  string      `json:"message,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	TotalPag int64       `json:"total_page,omitempty"`
}

// Failed use to response error messge
func Failed(w http.ResponseWriter, err error) {
	msg := err.Error()
	code := http.StatusInternalServerError

	switch t := err.(type) {
	case *exception.BadRequest:
		code = t.Code()
	case *exception.NotFound:
		code = t.Code()
	case *exception.InternalServerError:
		code = t.Code()
	case *exception.Unauthorized:
		code = t.Code()
	case *exception.MethodNotAllowed:
		code = t.Code()
	default:
		code = http.StatusInternalServerError
	}

	resp := Response{
		Status:  "error",
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

	w.WriteHeader(code)
	w.Write(respByt)
	return
}

// Success use to response success data
func Success(w http.ResponseWriter, code int, data interface{}) {
	resp := Response{
		Status:  "success",
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
		Status:   "success",
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

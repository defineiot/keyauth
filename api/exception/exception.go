package exception

import (
	"fmt"
	"net/http"
)

// APIException is openauth error
type APIException struct {
	msg  string
	code int
}

func (e *APIException) Error() string {
	return e.msg
}

// Code exception's code
func (e *APIException) Code() int {
	return e.code
}

// InternalServerError exception
type InternalServerError struct {
	*APIException
}

// NotFound exception
type NotFound struct {
	*APIException
}

// BadRequest exception
type BadRequest struct {
	*APIException
}

// NewAPIException is openauth api error
func NewAPIException(text string, code int) error {
	return &APIException{msg: text, code: code}
}

// NewInternalServerError for 503
func NewInternalServerError(format string, args ...interface{}) error {
	excp := new(InternalServerError)
	excp.msg = fmt.Sprintf(format, args...)
	excp.code = http.StatusInternalServerError
	return excp
}

// NewNotFound for 404
func NewNotFound(format string, args ...interface{}) error {
	excp := new(NotFound)
	excp.msg = fmt.Sprintf(format, args...)
	excp.code = http.StatusNotFound
	return excp
}

// NewBadRequest for 400
func NewBadRequest(format string, args ...interface{}) error {
	excp := new(BadRequest)
	excp.msg = fmt.Sprintf(format, args...)
	excp.code = http.StatusBadRequest
	return excp
}

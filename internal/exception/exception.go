package exception

import (
	"fmt"
	"net/http"
)

// APIException is openauth error
type apiException struct {
	msg  string
	code int
}

func (e *apiException) Error() string {
	return e.msg
}

// Code exception's code
func (e *apiException) Code() int {
	return e.code
}

// InternalServerError exception
type InternalServerError struct {
	*apiException
}

// NotFound exception
type NotFound struct {
	*apiException
}

// BadRequest exception
type BadRequest struct {
	*apiException
}

// Unauthorized exception
type Unauthorized struct {
	*apiException
}

// Forbidden exception
type Forbidden struct {
	*apiException
}

// Expired exception
type Expired struct {
	*apiException
}

// MethodNotAllowed exception
type MethodNotAllowed struct {
	*apiException
}

// NewInternalServerError for 503
func NewInternalServerError(format string, args ...interface{}) error {
	excp := new(InternalServerError)
	excp.apiException = &apiException{msg: fmt.Sprintf(format, args...), code: http.StatusInternalServerError}
	return excp
}

// NewNotFound for 404
func NewNotFound(format string, args ...interface{}) error {
	excp := new(NotFound)
	excp.apiException = &apiException{msg: fmt.Sprintf(format, args...), code: http.StatusNotFound}
	return excp
}

// NewBadRequest for 400
func NewBadRequest(format string, args ...interface{}) error {
	excp := new(BadRequest)
	excp.apiException = &apiException{msg: fmt.Sprintf(format, args...), code: http.StatusBadRequest}
	return excp
}

// NewUnauthorized for 401
func NewUnauthorized(format string, args ...interface{}) error {
	excp := new(Unauthorized)
	excp.apiException = &apiException{msg: fmt.Sprintf(format, args...), code: http.StatusUnauthorized}
	return excp
}

// NewExpired for 401
func NewExpired(format string, args ...interface{}) error {
	excp := new(Expired)
	excp.apiException = &apiException{msg: fmt.Sprintf(format, args...), code: http.StatusUnauthorized}
	return excp
}

// NewForbidden for 403
func NewForbidden(format string, args ...interface{}) error {
	excp := new(Forbidden)
	excp.apiException = &apiException{msg: fmt.Sprintf(format, args...), code: http.StatusForbidden}
	return excp
}

// NewMethodNotAllowed for 405
func NewMethodNotAllowed(format string, args ...interface{}) error {
	excp := new(MethodNotAllowed)
	excp.apiException = &apiException{msg: fmt.Sprintf(format, args...), code: http.StatusMethodNotAllowed}
	return excp
}

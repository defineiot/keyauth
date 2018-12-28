package exception

import (
	"fmt"
	"net/http"
)

// APIException is openauth error
type apiException struct {
	msg        string
	statusCode int
	errorCode  int
}

func (e *apiException) Error() string {
	return e.msg
}

// Code exception's code
func (e *apiException) Code() int {
	return e.statusCode
}

// InternalServerError exception
type InternalServerError struct {
	*apiException
}

// NewInternalServerError for 503
func NewInternalServerError(format string, args ...interface{}) error {
	excp := new(InternalServerError)
	excp.apiException = &apiException{msg: fmt.Sprintf(format, args...), statusCode: http.StatusInternalServerError}
	return excp
}

// NotFound exception
type NotFound struct {
	*apiException
}

// NewNotFound for 404
func NewNotFound(format string, args ...interface{}) error {
	excp := new(NotFound)
	excp.apiException = &apiException{msg: fmt.Sprintf(format, args...), statusCode: http.StatusNotFound}
	return excp
}

// HasExist exception
type HasExist struct {
	*apiException
}

// NewHasExist for 404
func NewHasExist(format string, args ...interface{}) error {
	excp := new(HasExist)
	excp.apiException = &apiException{msg: fmt.Sprintf(format, args...), statusCode: http.StatusBadRequest}
	return excp
}

// BadRequest exception
type BadRequest struct {
	*apiException
}

// NewBadRequest for 400
func NewBadRequest(format string, args ...interface{}) error {
	excp := new(BadRequest)
	excp.apiException = &apiException{msg: fmt.Sprintf(format, args...), statusCode: http.StatusBadRequest}
	return excp
}

// Unauthorized exception
type Unauthorized struct {
	*apiException
}

// NewUnauthorized for 401
func NewUnauthorized(format string, args ...interface{}) error {
	excp := new(Unauthorized)
	excp.apiException = &apiException{msg: fmt.Sprintf(format, args...), statusCode: http.StatusUnauthorized}
	return excp
}

// Forbidden exception
type Forbidden struct {
	*apiException
}

// NewForbidden for 403
func NewForbidden(format string, args ...interface{}) error {
	excp := new(Forbidden)
	excp.apiException = &apiException{msg: fmt.Sprintf(format, args...), statusCode: http.StatusForbidden}
	return excp
}

// Expired exception
type Expired struct {
	*apiException
}

// NewExpired for 401
func NewExpired(format string, args ...interface{}) error {
	excp := new(Expired)
	excp.apiException = &apiException{msg: fmt.Sprintf(format, args...), statusCode: http.StatusUnauthorized}
	return excp
}

// MethodNotAllowed exception
type MethodNotAllowed struct {
	*apiException
}

// NewMethodNotAllowed for 405
func NewMethodNotAllowed(format string, args ...interface{}) error {
	excp := new(MethodNotAllowed)
	excp.apiException = &apiException{msg: fmt.Sprintf(format, args...), statusCode: http.StatusMethodNotAllowed}
	return excp
}

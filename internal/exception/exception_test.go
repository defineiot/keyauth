package exception_test

import (
	"testing"

	"github.com/defineiot/keyauth/internal/exception"
)

func TestNewInternalServerError(t *testing.T) {
	err := exception.NewInternalServerError("internal error")
	val, ok := err.(*exception.InternalServerError)
	if !ok {
		t.Fatal("type assert failed, isn't InternalServerError type")
	}

	if val.Code() != 500 {
		t.Fatal("internal error code isn't 500")
	}
	if val.Error() != "internal error" {
		t.Fatal("message ins't internal error")
	}
}

func TestNewNotFound(t *testing.T) {
	err := exception.NewNotFound("not found error")
	val, ok := err.(*exception.NotFound)
	if !ok {
		t.Fatal("type assert failed, isn't NotFound type")
	}

	if val.Code() != 404 {
		t.Fatal("not found error code isn't 404")
	}
	if val.Error() != "not found error" {
		t.Fatal("message ins't not found error")
	}
}

func TestNewBadRequest(t *testing.T) {
	err := exception.NewBadRequest("bad request error")
	val, ok := err.(*exception.BadRequest)
	if !ok {
		t.Fatal("type assert failed, isn't BadRequest type")
	}

	if val.Code() != 400 {
		t.Fatal("bad request error code isn't 400")
	}
	if val.Error() != "bad request error" {
		t.Fatal("message ins't bad request error")
	}
}

func TestNewUnauthorized(t *testing.T) {
	err := exception.NewUnauthorized("unauthorized error")
	val, ok := err.(*exception.Unauthorized)
	if !ok {
		t.Fatal("type assert failed, isn't Unauthorized type")
	}

	if val.Code() != 401 {
		t.Fatal("internal error code isn't 401")
	}
	if val.Error() != "unauthorized error" {
		t.Fatal("message ins't unauthorized error")
	}
}

func TestNewMethodNotAllowed(t *testing.T) {
	err := exception.NewMethodNotAllowed("method not allowed error")
	val, ok := err.(*exception.MethodNotAllowed)
	if !ok {
		t.Fatal("type assert failed, isn't MethodNotAllowed type")
	}

	if val.Code() != 405 {
		t.Fatal("method not allowed error code isn't 405")
	}
	if val.Error() != "method not allowed error" {
		t.Fatal("message ins't method not allowed error")
	}
}

package oauth2

import "errors"

var (
	// OAuth2 errors in https://tools.ietf.org/html/rfc6749#section-4.1.2
	InvalidRequest          = errors.New("invalid_request")
	UnauthorizedClient      = errors.New("unauthorized_client")
	AccessDenied            = errors.New("access_denied")
	UnsupportedResponseType = errors.New("unsupported_response_type")
	InvalidScope            = errors.New("invalid_scope")
	ServerError             = errors.New("server_error")
	TemporarilyUnavailable  = errors.New("temporarily_unavailable")
)

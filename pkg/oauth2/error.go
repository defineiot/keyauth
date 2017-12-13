package oauth2

import "errors"

var (
	// OAuth2 errors in https://tools.ietf.org/html/rfc6749#section-4.1.2
	InvalidRequest          = errors.New("oauth2: invalid_request")
	UnauthorizedClient      = errors.New("oauth2: unauthorized_client")
	AccessDenied            = errors.New("oauth2:  access_denied")
	UnsupportedResponseType = errors.New("oauth2: unsupported_response_type")
	InvalidScope            = errors.New("oauth2: invalid_scope")
	ServerError             = errors.New("oauth2: server_error")
	TemporarilyUnavailable  = errors.New("oauth2: temporarily_unavailable")
)

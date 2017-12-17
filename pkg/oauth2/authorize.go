package oauth2

import "regexp"

// ResponseType the type of authorization request
type ResponseType string

const (
	// CODE define the type of authorization request
	CODE ResponseType = "code"
	// TOKEN define the type of authorization request
	TOKEN ResponseType = "token"

	// PKCE_PLAIN is oauth pkce extension
	PKCE_PLAIN = "plain"
	// PKCE_S256 is oauth pkce extension
	PKCE_S256 = "S256"
)

var (
	pkceMatcher = regexp.MustCompile("^[a-zA-Z0-9~._-]{43,128}$")
)

// AuthRequest Authorize request information
type AuthRequest struct {
}

// IssueCode use to issue auth code
func (c *Controller) IssueCode(*AuthRequest) (code string) {
	return
}

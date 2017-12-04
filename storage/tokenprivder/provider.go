package tokenprovider

import (
	"time"

	"openauth/storage/token"
)

// Manager use to provide token
type Manager interface {
	// Distribute a Bearer type Access Token
	// Bearer Token format:
	//    	Bearer XXXXXXXX
	// Which XXXXXXXX format b64token, ABNF definition:
	//    b64token = 1*( ALPHA / DIGIT / "-" / "." / "_" / "~" / "+" / "/" ) *"="
	// Written in Regular Expression:
	//    /[A-Za-z0-9\-\._~\+\/]+=*/
	//
	// Specify how long it will expire (in seconds) and
	// whether the refresh token is required, Token belongs to a user id
	IssueBearerAccessToken(expirationIn time.Duration, refresh bool, userID string) (*token.Token, error)
	// Revoke access token
	RevokeBearerToken(token string) error
	// Verify that Token is legal, Can only Validate one at a time,access or refresh
	ValidateBearerToken(access, refresh string) (*token.Token, error)
	// Undo all the Token for a entity, Returns the number of undo Token and error
	RevokenEntityBearerToken(entityID string) (int, error)
}

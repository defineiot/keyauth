package token

import "time"

// Token is use to access resource service
type Token struct {
	// Access Token string
	AccessToken string
	// Token type, default bearer
	TokenType string
	// token expires time
	ExpiresIn time.Time
	// a refresh token
	RefreshToken string
	// Token belongs to an entity, such as a user, a client
	AssociationID string
	// Extend fields to facilitate the expansion of database tables
	Extra string
}

// Manager is token service
type Manager interface {
	Issue(userID, projectID string) (*Token, error)
	IssueWithRefresh(refresh string) (*Token, error)
	IssueWithCode(code string) (*Token, error)
	RevokeToken(access string) error
	ValidateToken(access string) (*Token, error)
}

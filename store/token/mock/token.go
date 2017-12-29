package mock

import "openauth/store/token"

// TokenStore mock
type TokenStore struct {
	SaveTokenFn        func(t *token.Token) (*token.Token, error)
	SaveTokenInvoked   bool
	GetTokenFn         func(accessToken string) (*token.Token, error)
	GetTokenInvoked    bool
	DeleteTokenFn      func(accessToken string) error
	DeleteTokenInvoked bool
}

// Close mock
func (s *TokenStore) Close() error {
	return nil
}

// SaveToken mock
func (s *TokenStore) SaveToken(t *token.Token) (*token.Token, error) {
	s.SaveTokenInvoked = true
	retToken, err := s.SaveTokenFn(t)
	return retToken, err
}

// GetToken mock
func (s *TokenStore) GetToken(accessToken string) (*token.Token, error) {
	s.GetTokenInvoked = true
	t, err := s.GetTokenFn(accessToken)
	return t, err
}

// DeleteToken mock
func (s *TokenStore) DeleteToken(accessToken string) error {
	s.DeleteTokenInvoked = true
	err := s.DeleteTokenFn(accessToken)
	return err
}

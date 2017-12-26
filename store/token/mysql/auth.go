package mysql

import (
	"database/sql"
	"openauth/api/exception"
	"openauth/store/token"
)

func (s *store) SaveToken(t *token.Token) (*token.Token, error) {
	_, err := s.stmts[SaveToken].Exec(string(t.GrantType), t.AccessToken, t.RefreshToken, t.TokenType, t.CreatedAt, t.ExpiresIn, t.ClientID, t.UserID, t.Scope.DomainID, t.Scope.ProjectID)
	if err != nil {
		return nil, exception.NewInternalServerError("insert token exec sql err, %s", err)
	}

	return t, nil
}

func (s *store) GetToken(accessToken string) (*token.Token, error) {
	t := new(token.Token)
	t.Scope = new(token.Scope)

	err := s.stmts[FindToken].QueryRow("SELECT grant_type,access_token,refresh_token,type,create_at,expire_at,client_id,user_id,domain_id,project_id FROM token WHERE access_token = ?", accessToken).Scan(
		&t.GrantType, &t.AccessToken, &t.RefreshToken, &t.TokenType, &t.CreatedAt, &t.ExpiresIn, &t.ClientID, &t.UserID, &t.Scope.DomainID, &t.Scope.ProjectID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewUnauthorized("token %s not find", accessToken)
		}

		return nil, exception.NewInternalServerError("query token error, %s", err)
	}

	return t, nil
}

package mysql

import (
	"database/sql"

	"github.com/defineiot/keyauth/dao/token"
	"github.com/defineiot/keyauth/internal/exception"
)

func (s *store) SaveToken(t *token.Token) (*token.Token, error) {
	if t.Scope == nil {
		t.Scope = new(token.Scope)
	}

	_, err := s.stmts[SaveToken].Exec(string(t.GrantType), t.AccessToken, t.RefreshToken, t.TokenType, t.CreatedAt, t.ExpiresIn, t.ClientID, t.UserID, t.DomainID, t.Scope.WorkProject, t.ServiceName)
	if err != nil {
		return nil, exception.NewInternalServerError("insert token exec sql err, %s", err)
	}

	return t, nil
}

func (s *store) DeleteTokenByRefresh(refreshToken string) error {
	ret, err := s.stmts[DeleteTokenByRefresh].Exec(refreshToken)
	if err != nil {
		return exception.NewInternalServerError("delete refresh token exec sql error, %s", err)
	}
	count, err := ret.RowsAffected()
	if err != nil {
		return exception.NewInternalServerError("get delete refresh token row affected error, %s", err)
	}
	if count == 0 {
		return exception.NewBadRequest("refresh token %s not exist", refreshToken)
	}

	return nil
}

func (s *store) GetTokenByRefresh(refreshToken string) (*token.Token, error) {
	t := new(token.Token)
	t.Scope = new(token.Scope)

	err := s.stmts[FindTokenByRefresh].QueryRow(refreshToken).Scan(
		&t.GrantType, &t.AccessToken, &t.RefreshToken, &t.TokenType, &t.CreatedAt, &t.ExpiresIn, &t.ClientID, &t.UserID, &t.DomainID, &t.Scope.WorkProject, &t.ServiceName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewUnauthorized("refresh token %s not find", refreshToken)
		}

		return nil, exception.NewInternalServerError("query refresh token error, %s", err)
	}

	return t, nil
}

func (s *store) SetTokenProject(accessToken, projectID string) error {
	ret, err := s.stmts[SetTokenProject].Exec(projectID, accessToken)
	if err != nil {
		return exception.NewInternalServerError("update token project exec sql error, %s", err)
	}
	count, err := ret.RowsAffected()
	if err != nil {
		return exception.NewInternalServerError("get update token project row affected error, %s", err)
	}
	if count == 0 {
		return exception.NewBadRequest("access token %s not exist", accessToken)
	}

	return nil
}

func (s *store) GetToken(accessToken string) (*token.Token, error) {
	t := new(token.Token)
	t.Scope = new(token.Scope)

	err := s.stmts[FindToken].QueryRow(accessToken).Scan(
		&t.GrantType, &t.AccessToken, &t.RefreshToken, &t.TokenType, &t.CreatedAt, &t.ExpiresIn, &t.ClientID, &t.UserID, &t.DomainID, &t.Scope.WorkProject, &t.ServiceName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewUnauthorized("token %s not find", accessToken)
		}

		return nil, exception.NewInternalServerError("query token error, %s", err)
	}

	return t, nil
}

func (s *store) DeleteToken(accessToken string) error {
	ret, err := s.stmts[DeleteToken].Exec(accessToken)
	if err != nil {
		return exception.NewInternalServerError("delete token exec sql error, %s", err)
	}
	count, err := ret.RowsAffected()
	if err != nil {
		return exception.NewInternalServerError("get delete row affected error, %s", err)
	}
	if count == 0 {
		return exception.NewBadRequest("token %s not exist", accessToken)
	}

	return nil
}

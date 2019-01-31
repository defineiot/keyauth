package mysql

import (
	"database/sql"
	"time"

	"github.com/defineiot/keyauth/dao/token"
	"github.com/defineiot/keyauth/internal/exception"
)

func (s *store) SaveToken(t *token.Token) error {
	if err := t.ValidateSave(); err != nil {
		return err
	}

	if t.CreatedAt == 0 {
		t.CreatedAt = time.Now().Unix()
	}
	if t.ExpiresIn == 0 {
		t.ExpiresIn = 3600
	}

	if _, err := s.stmts[SaveToken].Exec(t.AccessToken, t.RefreshToken, string(t.GrantType),
		string(t.TokenType), t.UserID, t.DomainID, t.CurrentProject, t.ServiceID,
		t.ApplicationID, t.Name, t.Scope, t.CreatedAt, t.ExpiresIn, t.Description); err != nil {
		return exception.NewInternalServerError("insert token exec sql err, %s", err)
	}

	return nil
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

	if err := s.stmts[FindTokenByRefresh].QueryRow(refreshToken).Scan(
		&t.AccessToken, &t.RefreshToken, &t.GrantType, &t.TokenType,
		&t.UserID, &t.DomainID, &t.CurrentProject, &t.ServiceID, &t.ApplicationID,
		&t.Name, &t.Scope, &t.CreatedAt, &t.ExpiresIn, &t.Description); err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("refresh token %s not find", refreshToken)
		}

		return nil, exception.NewInternalServerError("query refresh token error, %s", err)
	}

	return t, nil
}

func (s *store) GetUserCurrentToken(userID, appID string, gt token.GrantType) (*token.Token, error) {
	t := new(token.Token)

	if err := s.stmts[FindUserCurrentToken].QueryRow(userID, appID, string(gt)).Scan(
		&t.AccessToken, &t.RefreshToken, &t.GrantType, &t.TokenType,
		&t.UserID, &t.DomainID, &t.CurrentProject, &t.ServiceID, &t.ApplicationID,
		&t.Name, &t.Scope, &t.CreatedAt, &t.ExpiresIn, &t.Description); err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("user:%s, app: %s, gt: %s token not find",
				userID, appID, gt)
		}

		return nil, exception.NewInternalServerError("query user current token error, %s", err)
	}

	return t, nil
}

func (s *store) UpdateTokenScope(accessToken, scope string) error {
	ret, err := s.stmts[UpdateTokenScope].Exec(scope, accessToken)
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

	if err := s.stmts[FindToken].QueryRow(accessToken).Scan(
		&t.AccessToken, &t.RefreshToken, &t.GrantType, &t.TokenType,
		&t.UserID, &t.DomainID, &t.CurrentProject, &t.ServiceID, &t.ApplicationID,
		&t.Name, &t.Scope, &t.CreatedAt, &t.ExpiresIn, &t.Description); err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewUnauthorized("access token %s not find", accessToken)
		}

		return nil, exception.NewInternalServerError("query access token error, %s", err)
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

package mysql

import (
	"database/sql"
	"openauth/api/exception"
	"openauth/storage/token"
	"sync"
)

var (
	createStmt *sql.Stmt
	deleteStmt *sql.Stmt
)

// NewTokenStorage use to new token storage
func NewTokenStorage(db *sql.DB) token.Storage {
	return &manager{db: db}
}

type manager struct {
	db *sql.DB
}

func (m *manager) SaveToken(t *token.Token) (*token.Token, error) {
	if err := t.Validate(); err != nil {
		return nil, err
	}

	var (
		once sync.Once
		err  error
	)

	once.Do(func() {
		createStmt, err = m.db.Prepare("INSERT INTO `token` (grant_type, access_token, refresh_token, type, create_at, expire_at, client_id, user_id, domian_id, project_id) VALUES (?,?,?,?,?,?,?,?,?,?)")
	})
	if err != nil {
		return nil, exception.NewInternalServerError("prepare insert token stmt error, %s", err)
	}

	_, err = createStmt.Exec(t.GrantType, t.AccessToken, t.RefreshToken, t.TokenType, t.CreatedAt, t.ExpiresIn, t.ClientID, t.UserID, t.Scope.DomainID, t.Scope.ProjectID)
	if err != nil {
		return nil, exception.NewInternalServerError("insert token exec sql err, %s", err)
	}

	return t, nil
}

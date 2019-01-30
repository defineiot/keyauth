package mysql

import (
	"database/sql"

	"github.com/defineiot/keyauth/dao"
	"github.com/defineiot/keyauth/dao/token"
	"github.com/defineiot/keyauth/internal/exception"
	"github.com/defineiot/keyauth/internal/tools"
)

const (
	SaveToken            = "save-token"
	FindToken            = "find-token"
	DeleteToken          = "delete-token"
	DeleteTokenByRefresh = "delete-token-by-refresh"
	FindTokenByRefresh   = "find-token-by-refresh"
	UpdateTokenScope     = "set-token-project"
	FindUserCurrentToken = "find-user-current-token"
)

// NewTokenStore use to create domain storage service
func NewTokenStore(opt *dao.Options) (token.Store, error) {
	unprepared := map[string]string{
		SaveToken: `
			INSERT INTO tokens (access_token, refresh_token, grant_type, token_type, user_id, domain_id, project_id, service_id, application_id, name, scope, create_at, expire_at, description) 
			VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?);
		`,
		FindToken: `
			SELECT access_token, refresh_token, grant_type, token_type, user_id, domain_id, project_id, service_id, application_id, name, scope, create_at, expire_at, description 
			FROM tokens 
			WHERE access_token = ?;
		`,
		FindUserCurrentToken: `
			SELECT access_token, refresh_token, grant_type, token_type, user_id, domain_id, project_id, service_id, application_id, name, scope, create_at, expire_at, description 
			FROM tokens 
			WHERE user_id = ? 
			AND application_id = ? 
			AND grant_type = ? 
			ORDER BY create_at 
			DESC 
			LIMIT 1;
		`,
		DeleteToken: `
			DELETE FROM tokens 
			WHERE access_token = ?;
		`,
		DeleteTokenByRefresh: `
		    DELETE FROM tokens 
		    WHERE refresh_token = ?;
		`,
		FindTokenByRefresh: `
		    SELECT access_token, refresh_token, grant_type, token_type, user_id, domain_id, project_id, service_id, application_id, name, scope, create_at, expire_at, description 
		    FROM tokens 
		    WHERE refresh_token = ?;
		`,
		UpdateTokenScope: `
			UPDATE tokens 
			SET scope = ? 
			WHERE access_token = ?;
		`,
	}

	// prepare all statements to verify syntax
	stmts, err := tools.PrepareStmts(opt.DB, unprepared)
	if err != nil {
		return nil, exception.NewInternalServerError("prepare token store query statment error, %s", err)
	}

	s := store{
		db:    opt.DB,
		stmts: stmts,
	}

	return &s, nil
}

// DomainManager is use mongodb as storage
type store struct {
	db    *sql.DB
	stmts map[string]*sql.Stmt
}

// Close closes the database, releasing any open resources.
func (s *store) Close() error {
	return s.db.Close()
}

func init() {
	dao.Registe(NewTokenStore)
}

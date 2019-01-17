package mysql

import (
	"database/sql"

	"github.com/defineiot/keyauth/dao"
	"github.com/defineiot/keyauth/dao/application"
	"github.com/defineiot/keyauth/internal/exception"
	"github.com/defineiot/keyauth/internal/tools"
)

const (
	CreateAPP    = "create-application"
	DeleteAPP    = "delete-app"
	ListUserAPPS = "list-user-apps"
	GetUserAPP   = "get-user-app"

	CheckExistByID   = "check-exist-by-id"
	CheckExistByName = "check-exist-by-name"
)

// NewAppStore use to create domain storage service
func NewAppStore(opt *dao.Options) (application.Store, error) {
	if opt.DB == nil {
		return nil, exception.NewInternalServerError("mysql connection is null")
	}

	unprepared := map[string]string{
		CreateAPP: `
			INSERT INTO applications (id, name, user_id, website, logo_image, description, create_at, redirect_uri, client_id, client_secret, locked, token_expire_time) 
			VALUES (?,?,?,?,?,?,?,?,?,?,?,?);
		`,
		ListUserAPPS: `
			SELECT id, name, user_id, website, logo_image, description, create_at, redirect_uri, client_id, client_secret, locked, last_login_time, last_login_ip, login_failed_times, login_success_times, token_expire_time
			FROM applications 
			WHERE user_id = ? 
			ORDER BY create_at 
			DESC;
		`,
		GetUserAPP: `
		    SELECT id, name, user_id, website, logo_image, description, create_at, redirect_uri, client_id, client_secret, locked, last_login_time, last_login_ip, login_failed_times, login_success_times, token_expire_time
		    FROM applications 
		    WHERE id = ? 
		    ORDER BY create_at 
		    DESC;
	    `,
		DeleteAPP: `
			DELETE FROM applications 
			WHERE id = ?;
		`,
		CheckExistByID: `
			SELECT id FROM applications 
			WHERE id = ?;
		`,
		CheckExistByName: `
			SELECT id FROM applications 
			WHERE name = ? 
			AND user_id = ?;
		`,
	}

	// prepare all statements to verify syntax
	stmts, err := tools.PrepareStmts(opt.DB, unprepared)
	if err != nil {
		return nil, exception.NewInternalServerError("prepare application store query statment error, %s", err)
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
	dao.Registe(NewAppStore)
}

package mysql

import (
	"database/sql"
	"fmt"

	"openauth/api/exception"
	"openauth/api/logger"
	"openauth/store/user"
)

const (
	SaveUser         = "save-user"
	FindAllUsers     = "find-all-users"
	FindUserByID     = "find-user-by-id"
	FindUserByName   = "find-user-by-name"
	FindUserPhones   = "find-user-phones"
	FindUserEmails   = "find-user-emails"
	FindUserPassword = "find-user-password"
	DeleteUserByID   = "delete-user-by-id"
	FindUserIDByName = "find-user-id-by-name"

	FindUserProjects       = "find-user-projects"
	SetUserDefaultProject  = "set-user-default-project"
	AddProjectToUser       = "add-project-to-user"
	RemoveProjectsFromUser = "remove-projects-from-user"

	CheckUserExistByName = "check-user-exist-by-name"
	CheckUserExistByID   = "check-user-exist-by-id"
)

// NewUserStore use to create domain storage service
func NewUserStore(db *sql.DB, key string, logger logger.OpenAuthLogger) (user.Store, error) {
	unprepared := map[string]string{
		SaveUser: `
			INSERT INTO users (id, name, enabled, domain_id, create_at, expires_active_days) 
			VALUES (?,?,?,?,?,?);
		`,
		FindAllUsers: `
			SELECT u.id, u.name, u.enabled, u.last_active_time, u.create_at, u.expires_active_days, u.default_project_id 
			FROM users u
			WHERE domain_id = ?;
		`,
		FindUserByID: `
			SELECT u.id, u.name, u.enabled, u.last_active_time, u.domain_id, u.create_at, u.expires_active_days, u.default_project_id 
			FROM users u
			WHERE id = ?;
		`,
		FindUserByName: `
			SELECT u.id, u.name, u.enabled, u.last_active_time, u.domain_id, u.create_at, u.expires_active_days, u.default_project_id 
			FROM users u 
			WHERE name = ? 
			AND domain_id = ?;
		`,
		FindUserProjects: `
			SELECT project_id 
			FROM users_projects_mapping 
			WHERE user_id = ?;
		`,
		SetUserDefaultProject: `
			UPDATE users
			SET default_project_id = ? 
			WHERE id = ?;
		`,
		AddProjectToUser: `
			INSERT INTO users_projects_mapping (user_id, project_id) 
			VALUES (?,?);
		`,
		RemoveProjectsFromUser: `
			DELETE FROM users_projects_mapping 
			WHERE user_id = ? 
			AND project_id = ?;
		`,
		FindUserPhones: `
			SELECT p.id, p.numbers, 'p.primary', p.description 
			FROM phones p 
			WHERE user_id = ?;
		`,
		FindUserEmails: `
			SELECT e.id, e.address, 'e.primary', e.description 
			FROM emails e
			WHERE user_id = ?;
		`,
		FindUserPassword: `
			SELECT p.password, p.expires_at, p.create_at, p.update_at 
			FROM passwords p
			WHERE user_id = ?;
		`,
		FindUserIDByName: `
			SELECT u.id 
			FROM users u
			WHERE name = ? 
			AND domain_id = ?;
		`,
		CheckUserExistByName: `
			SELECT u.name 
			FROM users u
			WHERE name = ? 
			AND domain_id = ?;
		`,
		CheckUserExistByID: `
			SELECT u.id 
			FROM users u
			WHERE id = ?;
		`,
	}

	// prepare all statements to verify syntax
	stmts, err := prepareStmts(db, unprepared)
	if err != nil {
		return nil, exception.NewInternalServerError("prepare user store query statment error, %s", err)
	}

	s := store{
		db:    db,
		stmts: stmts,
		key:   key,
		log:   logger,
	}

	return &s, nil
}

// DomainManager is use mongodb as storage
type store struct {
	db    *sql.DB
	stmts map[string]*sql.Stmt
	key   string
	log   logger.OpenAuthLogger
}

// Close closes the database, releasing any open resources.
func (s *store) Close() error {
	return s.db.Close()
}

// prepareStmts will attempt to prepare each unprepared
// query on the database. If one fails, the function returns
// with an error.
func prepareStmts(db *sql.DB, unprepared map[string]string) (map[string]*sql.Stmt, error) {
	prepared := map[string]*sql.Stmt{}
	for k, v := range unprepared {
		stmt, err := db.Prepare(v)
		if err != nil {
			return nil, fmt.Errorf("prepare statment: %s, %s", k, err)
		}
		prepared[k] = stmt
	}

	return prepared, nil
}

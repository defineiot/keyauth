package mysql

import (
	"database/sql"

	"openauth/api/exception"
	"openauth/store/project"
)

const (
	CreateProject      = "create-project"
	FindProjectByID    = "find-project-by-id"
	UpdateProjectByID  = "update-project-by-id"
	FindDomainPorjects = "find-domain-projects"
	DeleteProject      = "delete-project"

	FindProjectUsers       = "find-project-users"
	AddUsersToProject      = "add-users-to-project"
	RemoveUsersFromProject = "remove-users-from-project"

	CheckProjectExistByID   = "check-project-exist-by-id"
	CheckProjectExistByName = "check-project-exist-by-name"
)

// NewProjectStore use to create domain storage service
func NewProjectStore(db *sql.DB) (project.Store, error) {
	unprepared := map[string]string{
		CreateProject: `
			INSERT INTO project (id, name, description, enabled, domain_id, create_at) 
			VALUES (?,?,?,?,?,?);
		`,
		FindDomainPorjects: `
			SELECT p.id, p.name, p.description, p.enabled, p.domain_id, p.create_at 
			FROM project p
			WHERE domain_id = ? 
			ORDER BY create_at 
			DESC;
		`,
		FindProjectByID: `
			SELECT p.id, p.name, p.description, p.enabled, p.create_at, p.domain_id 
			FROM project p
			WHERE id = ?;
		`,
		DeleteProject: `
			DELETE FROM project 
			WHERE id = ?;
		`,
		CheckProjectExistByID: `
			SELECT id FROM project 
			WHERE id = ?;
		`,
		CheckProjectExistByName: `
			SELECT name 
			FROM project 
			WHERE name = ? 
			AND domain_id = ?;
		`,
		FindProjectUsers: `
			SELECT user_id 
			FROM mapping 
			WHERE project_id = ?;
		`,
		AddUsersToProject: `
			INSERT INTO mapping (user_id, project_id) 
			VALUES (?,?);
		`,
		RemoveUsersFromProject: `
			DELETE FROM mapping 
			WHERE user_id = ? 
			AND project_id = ?;
		`,
	}

	// prepare all statements to verify syntax
	stmts, err := prepareStmts(db, unprepared)
	if err != nil {
		return nil, exception.NewInternalServerError("prepare domain query statment error, %s", err)
	}

	s := store{
		db:    db,
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

// prepareStmts will attempt to prepare each unprepared
// query on the database. If one fails, the function returns
// with an error.
func prepareStmts(db *sql.DB, unprepared map[string]string) (map[string]*sql.Stmt, error) {
	prepared := map[string]*sql.Stmt{}
	for k, v := range unprepared {
		stmt, err := db.Prepare(v)
		if err != nil {
			return nil, err
		}
		prepared[k] = stmt
	}

	return prepared, nil
}

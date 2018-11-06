package mysql

import (
	"database/sql"

	"github.com/defineiot/keyauth/dao/project"
	"github.com/defineiot/keyauth/internal/exception"
	"github.com/defineiot/keyauth/internal/tools"
)

const (
	CreateProject       = "create-project"
	FindProjectByID     = "find-project-by-id"
	UpdateProjectByID   = "update-project-by-id"
	FindDomainPorjects  = "find-domain-projects"
	DeleteProjectByID   = "delete-project-by-id"
	DeleteProjectByName = "delete-project-by-name"

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
			INSERT INTO projects (id, name, description, enabled, domain_id, create_at) 
			VALUES (?,?,?,?,?,?);
		`,
		FindDomainPorjects: `
			SELECT p.id, p.name, p.description, p.enabled, p.domain_id, p.create_at 
			FROM projects p
			WHERE domain_id = ? 
			ORDER BY create_at 
			DESC;
		`,
		FindProjectByID: `
			SELECT p.id, p.name, p.description, p.enabled, p.create_at, p.domain_id 
			FROM projects p
			WHERE id = ?;
		`,
		DeleteProjectByID: `
			DELETE FROM projects 
			WHERE id = ?;
		`,
		DeleteProjectByName: `
			DELETE FROM projects 
			WHERE name = ? 
			AND domain_id = ?;
		`,
		CheckProjectExistByID: `
			SELECT id FROM projects 
			WHERE id = ?;
		`,
		CheckProjectExistByName: `
			SELECT name 
			FROM projects 
			WHERE name = ? 
			AND domain_id = ?;
		`,
		FindProjectUsers: `
			SELECT user_id 
			FROM users_projects_mapping 
			WHERE project_id = ?;
		`,
		AddUsersToProject: `
			INSERT INTO users_projects_mapping (user_id, project_id) 
			VALUES (?,?);
		`,
		RemoveUsersFromProject: `
			DELETE FROM users_projects_mapping 
			WHERE user_id = ? 
			AND project_id = ?;
		`,
	}

	// prepare all statements to verify syntax
	stmts, err := tools.PrepareStmts(db, unprepared)
	if err != nil {
		return nil, exception.NewInternalServerError("prepare  project store query statment error, %s", err)
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

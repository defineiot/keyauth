package mysql

import (
	"database/sql"

	"github.com/defineiot/keyauth/dao"
	"github.com/defineiot/keyauth/dao/project"
	"github.com/defineiot/keyauth/internal/exception"
	"github.com/defineiot/keyauth/internal/tools"
)

const (
	CreateProject          = "create-project"
	FindProjectByID        = "find-project-by-id"
	UpdateProjectByID      = "update-project-by-id"
	FindDomainPorjects     = "find-domain-projects"
	FindDepartmentProjects = "find-department-projects"
	FindUserProjects       = "find-user-projects"
	DeleteProjectByID      = "delete-project-by-id"
	DeleteProjectByName    = "delete-project-by-name"

	FindProjectUsers        = "find-project-users"
	AddUsersToProject       = "add-users-to-project"
	RemoveUsersFromProject  = "remove-users-from-project"
	CheckProjectExistByID   = "check-project-exist-by-id"
	CheckProjectExistByName = "check-project-exist-by-name"
)

// NewProjectStore use to create domain storage service
func NewProjectStore(opt *dao.Options) (project.Store, error) {
	unprepared := map[string]string{
		CreateProject: `
			INSERT INTO projects (id, name, picture, latitude, longitude, enabled, owner_id,  description, domain_id, create_at) 
			VALUES (?,?,?,?,?,?,?,?,?,?);
		`,
		FindDomainPorjects: `
			SELECT id, name, picture, latitude, longitude, enabled, owner_id,  description, domain_id, create_at, update_at 
			FROM projects
			WHERE domain_id = ? 
			ORDER BY create_at 
			DESC;
		`,
		FindDepartmentProjects: `
			SELECT id, name, picture, latitude, longitude, enabled, owner_id,  description, domain_id, create_at, update_at 
			FROM projects p
			LEFT JOIN department_project_mappings m 
			ON m.project_id = p.id 
			WHERE m.department_id = ? 
			ORDER BY create_at 
			DESC;
		`,
		FindUserProjects: `
			SELECT id, name, picture, latitude, longitude, enabled, owner_id, description, domain_id, create_at, update_at 
			FROM projects p 
			LEFT JOIN user_project_mappings m 
			ON p.id = m.project_id
			WHERE p.domain_id = ? 
			AND m.user_id = ? 
			ORDER BY p.create_at 
			DESC;
		`,
		FindProjectByID: `
			SELECT id, name, picture, latitude, longitude, enabled, owner_id,  description, domain_id, create_at, update_at 
			FROM projects
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
			FROM user_project_mappings
			WHERE project_id = ?;
		`,
		AddUsersToProject: `
			INSERT INTO user_project_mappings (user_id, project_id) 
			VALUES (?,?);
		`,
		RemoveUsersFromProject: `
			DELETE FROM user_project_mappings 
			WHERE user_id = ? 
			AND project_id = ?;
		`,
	}

	// prepare all statements to verify syntax
	stmts, err := tools.PrepareStmts(opt.DB, unprepared)
	if err != nil {
		return nil, exception.NewInternalServerError("prepare  project store query statment error, %s", err)
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
	dao.Registe(NewProjectStore)
}

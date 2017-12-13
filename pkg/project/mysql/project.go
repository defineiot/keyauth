package mysql

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/satori/go.uuid"

	"openauth/api/exception"
	"openauth/pkg/project"
)

var (
	createPrepare *sql.Stmt
	deletePrepare *sql.Stmt
)

// NewProjectService is use mysql as storage
func NewProjectService(db *sql.DB) project.Service {
	return &manager{db: db}
}

type manager struct {
	db *sql.DB
}

func (m *manager) CreateProject(domainID, name, description string, enabled bool) (*project.Project, error) {
	var (
		once   sync.Once
		preErr error
	)

	ok, err := m.projectNameExist(domainID, name)
	if err != nil {
		return nil, exception.NewInternalServerError("check project name exist error, %s", err)
	}
	if ok {
		return nil, exception.NewBadRequest("project name %s in this domain is exists", name)
	}

	once.Do(func() {
		createPrepare, preErr = m.db.Prepare("INSERT INTO `project` (id, name, description, enabled, domain_id, create_at) VALUES (?,?,?,?,?,?)")
	})
	if preErr != nil {
		return nil, exception.NewInternalServerError("prepare insert project stmt error, project: %s, %s", name, err)
	}

	proj := project.Project{ID: uuid.NewV4().String(), Name: name, Description: description, CreateAt: time.Now().Unix(), Enabled: enabled, DomainID: domainID}
	_, err = createPrepare.Exec(proj.ID, proj.Name, proj.Description, proj.Enabled, proj.DomainID, proj.CreateAt)
	if err != nil {
		return nil, exception.NewInternalServerError("insert project exec sql err, %s", err)
	}

	return &proj, nil
}

// Notice: if project not exits return nil
func (m *manager) GetProject(id string) (*project.Project, error) {
	proj := project.Project{}
	err := m.db.QueryRow("SELECT id,name,description,enabled,create_at,domain_id FROM project WHERE id = ?", id).Scan(
		&proj.ID, &proj.Name, &proj.Description, &proj.Enabled, &proj.CreateAt, &proj.DomainID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("project %s not find", id)
		}

		return nil, exception.NewInternalServerError("query single project error, %s", err)
	}

	return &proj, nil
}

func (m *manager) ListDomainProjects(domainID string) ([]*project.Project, error) {
	rows, err := m.db.Query("SELECT id,name,description,enabled,domain_id,create_at FROM project WHERE domain_id = ? ORDER BY create_at DESC", domainID)
	if err != nil {
		return nil, exception.NewInternalServerError("query project list error, %s", err)
	}
	defer rows.Close()

	projects := []*project.Project{}
	for rows.Next() {
		proj := project.Project{}
		if err := rows.Scan(&proj.ID, &proj.Name, &proj.Description, &proj.Enabled, &proj.DomainID, &proj.CreateAt); err != nil {
			return nil, exception.NewInternalServerError("scan project record error, %s", err)
		}
		projects = append(projects, &proj)
	}

	return projects, nil
}

func (m *manager) UpdateProject(id, name, description string) (*project.Project, error) {
	return nil, nil
}

func (m *manager) DeleteProject(id string) error {
	var (
		once sync.Once
		err  error
	)

	once.Do(func() {
		deletePrepare, err = m.db.Prepare("DELETE FROM project WHERE id = ?")
	})
	if err != nil {
		return exception.NewInternalServerError("prepare delete project stmt error, %s", err)
	}

	ret, err := deletePrepare.Exec(id)
	if err != nil {
		return exception.NewInternalServerError("delete project exec sql error, %s", err)
	}
	count, err := ret.RowsAffected()
	if err != nil {
		return exception.NewInternalServerError("get delete row affected error, %s", err)
	}
	if count == 0 {
		return exception.NewBadRequest("project %s not exist", id)
	}

	return nil
}

func (m *manager) CheckProjectIsExistByID(id string) (bool, error) {
	var pid string
	err := m.db.QueryRow("SELECT id FROM project WHERE id = ?", id).Scan(&pid)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, exception.NewInternalServerError("check project exist error, %s", err)
	}

	return true, nil
}

func (m *manager) projectNameExist(domainID, projectName string) (bool, error) {
	rows, err := m.db.Query("SELECT name FROM project WHERE name = ? AND domain_id = ?", projectName, domainID)
	if err != nil {
		return false, fmt.Errorf("query project name exist error, %s", err)
	}
	defer rows.Close()

	projects := []string{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return false, fmt.Errorf("scan project name exist record error, %s", err)
		}
		projects = append(projects, name)
	}
	if len(projects) != 0 {
		return true, nil
	}

	return false, nil
}

func (m *manager) ListProjectUsers(projectID string) ([]string, error) {
	// check the project is exist
	ok, err := m.CheckProjectIsExistByID(projectID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, exception.NewBadRequest("project %s not exist", projectID)
	}

	rows, err := m.db.Query("SELECT user_id FROM mapping WHERE project_id = ?", projectID)
	if err != nil {
		return nil, exception.NewInternalServerError("query project user id error, %s", err)
	}
	defer rows.Close()

	users := []string{}
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, exception.NewInternalServerError("scan project's user id error, %s", err)
		}
		users = append(users, userID)
	}
	return users, nil
}

func (m *manager) AddUsersToProject(projectID string, userIDs ...string) error {
	// check the project is exist
	ok, err := m.CheckProjectIsExistByID(projectID)
	if err != nil {
		return err
	}
	if !ok {
		return exception.NewBadRequest("project %s not exist", projectID)
	}

	// check user is in this project
	uids, err := m.ListProjectUsers(projectID)
	if err != nil {
		return err
	}
	existUIDs := []string{}
	for _, uid := range uids {
		for _, inuid := range userIDs {
			if inuid == uid {
				existUIDs = append(existUIDs, inuid)
			}
		}
	}
	if len(existUIDs) != 0 {
		return exception.NewBadRequest("users %s is in this project", existUIDs)
	}

	mappPre, err := m.db.Prepare("INSERT INTO `mapping` (user_id, project_id) VALUES (?,?)")
	if err != nil {
		return fmt.Errorf("prepare add users to project mapping stmt error, project: %s, user: %s,  %s", projectID, userIDs, err)
	}

	for _, userID := range userIDs {
		_, err = mappPre.Exec(userID, projectID)
		if err != nil {
			mappPre.Close()
			return fmt.Errorf("insert add users to project mapping exec sql err, %s", err)
		}
	}
	mappPre.Close()

	return nil
}

func (m *manager) RemoveUsersFromProject(projectID string, userIDs ...string) error {
	// check the project is exist
	ok, err := m.CheckProjectIsExistByID(projectID)
	if err != nil {
		return err
	}
	if !ok {
		return exception.NewBadRequest("project %s not exist", projectID)
	}

	// check user is in this project
	uids, err := m.ListProjectUsers(projectID)
	if err != nil {
		return err
	}
	nExistUIDs := []string{}
	for _, inuid := range userIDs {
		var ok bool
		for _, uid := range uids {
			if uid == inuid {
				ok = true
			}
		}
		if !ok {
			nExistUIDs = append(nExistUIDs, inuid)
		}
	}
	if len(nExistUIDs) != 0 {
		return exception.NewBadRequest("users %s isn't in this project", nExistUIDs)
	}

	mappPre, err := m.db.Prepare("DELETE FROM `mapping` WHERE user_id = ? AND project_id = ?")
	if err != nil {
		return fmt.Errorf("prepare remove users from project mapping stmt error, project: %s, user: %s,  %s", projectID, userIDs, err)
	}

	for _, userID := range userIDs {
		_, err = mappPre.Exec(userID, projectID)
		if err != nil {
			mappPre.Close()
			return fmt.Errorf("insert remove users from project mapping exec sql err, %s", err)
		}
	}
	mappPre.Close()

	return nil
}

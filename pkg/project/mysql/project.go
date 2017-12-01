package mysql

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"openauth/pkg/project"

	uuid "github.com/satori/go.uuid"
)

var (
	createPrepare *sql.Stmt
	deletePrepare *sql.Stmt
)

// NewProjectManager is use mysql as storage
func NewProjectManager(db *sql.DB) project.Manager {
	return &manager{db: db}
}

type manager struct {
	db *sql.DB
}

func (m *manager) CreateProject(domainID, name, description string, enabled bool) (*project.Project, error) {
	var (
		once sync.Once
		err  error
	)

	once.Do(func() {
		createPrepare, err = m.db.Prepare("INSERT INTO `project` (id, name, description, enabled, domain_id, create_at) VALUES (?,?,?,?,?,?)")
	})
	if err != nil {
		return nil, fmt.Errorf("prepare insert project stmt error, project: %s, %s", name, err)
	}

	proj := project.Project{ID: uuid.NewV4().String(), Name: name, Description: description, CreateAt: time.Now().Unix(), Enabled: enabled, DomainID: domainID}
	_, err = createPrepare.Exec(proj.ID, proj.Name, proj.Description, proj.Enabled, proj.DomainID, proj.CreateAt)
	if err != nil {
		return nil, fmt.Errorf("insert project exec sql err, %s", err)
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
			return nil, nil
		}

		return nil, fmt.Errorf("query single project error, %s", err)
	}

	return &proj, nil
}

func (m *manager) ListProject(domainID string) ([]*project.Project, error) {
	rows, err := m.db.Query("SELECT id,name,description,enabled,domain_id,create_at FROM project")
	if err != nil {
		return nil, fmt.Errorf("query project list error, %s", err)
	}
	defer rows.Close()

	projects := []*project.Project{}
	for rows.Next() {
		proj := project.Project{}
		if err := rows.Scan(&proj.ID, &proj.Name, &proj.Description, &proj.Enabled, &proj.DomainID, &proj.CreateAt); err != nil {
			return nil, fmt.Errorf("scan project record error, %s", err)
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
		return fmt.Errorf("prepare delete project stmt error, %s", err)
	}

	if _, err := deletePrepare.Exec(id); err != nil {
		return fmt.Errorf("delete project exec sql error, %s", err)
	}
	return nil
}

func (m *manager) IsExist(id string) (bool, error) {
	var pid string
	err := m.db.QueryRow("SELECT id FROM project WHERE id = ?", id).Scan(&pid)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, fmt.Errorf("check project exist error, %s", err)
	}

	return true, nil
}

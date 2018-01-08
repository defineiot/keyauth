package mysql

import (
	"openauth/api/exception"
	"time"

	"github.com/satori/go.uuid"

	"openauth/store/application"
	"openauth/store/service"
	"openauth/tools"
)

func (s *store) SaveService(name, description string) (*service.Service, error) {
	client := new(application.Client)
	client.ClientID = tools.MakeUUID(24)
	client.ClientSecret = tools.MakeUUID(32)
	client.ClientType = "confidential"

	svr := new(service.Service)
	svr.ID = uuid.NewV4().String()
	svr.Name = name
	svr.Description = description
	svr.CreateAt = time.Now().Unix()
	svr.Enabled = true
	svr.Client = client

	tx, err := s.db.Begin()
	if err != nil {
		return nil, exception.NewInternalServerError("start transaction error, %s", err)
	}
	defer tx.Rollback()

	// insert client
	stmtC, err := tx.Prepare("INSERT INTO clients (id, secret, type, redirect_uri, application_id, service_id) VALUES (?,?,?,?,?,?)")
	if err != nil {
		return nil, exception.NewInternalServerError("prepare insert client error, %s", err)
	}
	if _, err := stmtC.Exec(client.ClientID, client.ClientSecret, client.ClientType, client.RedirectURI, "", svr.ID); err != nil {
		stmtC.Close()
		return nil, exception.NewInternalServerError("exec sql error, %s", err)
	}
	stmtC.Close()

	// insert service
	stmtS, err := tx.Prepare("INSERT INTO services (id, name, description, create_at, enabled) VALUES (?,?,?,?,?)")
	if err != nil {
		return nil, exception.NewInternalServerError("prepare insert service error, %s", err)
	}
	if _, err := stmtS.Exec(svr.ID, svr.Name, svr.Description, svr.CreateAt, svr.Enabled); err != nil {
		stmtS.Close()
		return nil, exception.NewInternalServerError("exec sql error, %s", err)
	}
	stmtS.Close()

	if err := tx.Commit(); err != nil {
		return nil, exception.NewInternalServerError("commit error, %s", err)
	}

	return svr, nil
}

func (s *store) DeleteService(sid string) error {
	return nil
}

func (s *store) FindAllService() ([]*service.Service, error) {
	return nil, nil
}

func (s *store) FindServiceByID() (*service.Service, error) {
	return nil, nil
}

package mysql

import (
	"database/sql"
	"time"

	"github.com/satori/go.uuid"

	"openauth/api/exception"
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
	tx, err := s.db.Begin()
	if err != nil {
		return exception.NewInternalServerError("start transaction error, %s", err)
	}
	defer tx.Rollback()

	// delete service
	stmtS, err := tx.Prepare("DELETE FROM services WHERE id = ?")
	if err != nil {
		return exception.NewInternalServerError("parpare delete service sql error, %s", err)
	}
	res, err := stmtS.Exec(sid)
	if err != nil {
		stmtS.Close()
		return exception.NewInternalServerError("exec delete sql error, %s", err)
	}
	aff, err := res.RowsAffected()
	if err != nil {
		stmtS.Close()
		return exception.NewInternalServerError("exec delete sql affect error, %s", err)
	}
	if aff == 0 {
		stmtS.Close()
		return exception.NewBadRequest("service %s not found", sid)
	}
	stmtS.Close()

	// delete client
	stmtC, err := tx.Prepare("DELETE FROM clients WHERE service_id = ?")
	if err != nil {
		return exception.NewInternalServerError("parpare delete service client error, %s", err)
	}
	if _, err := stmtC.Exec(sid); err != nil {
		stmtC.Close()
		return exception.NewInternalServerError("delete service client error, %s", err)
	}
	stmtC.Close()

	if err := tx.Commit(); err != nil {
		return exception.NewInternalServerError("commit error, %s", err)
	}

	return nil
}

func (s *store) FindAllService() ([]*service.Service, error) {
	rows, err := s.stmts[FindAll].Query()
	if err != nil {
		return nil, exception.NewInternalServerError("query service list error, %s", err)
	}
	defer rows.Close()

	svrs := []*service.Service{}
	for rows.Next() {
		svr := service.Service{}
		cli := application.Client{}
		if err := rows.Scan(&svr.ID, &svr.Name, &svr.Description, &svr.Enabled, &svr.Status, &svr.StatusUpdateAt, &svr.Version, &svr.CreateAt, &cli.ClientID, &cli.ClientSecret); err != nil {
			return nil, exception.NewInternalServerError("scan service record error, %s", err)
		}
		svr.Client = &cli
		svrs = append(svrs, &svr)
	}

	return svrs, nil
}

func (s *store) FindServiceByID(sid string) (*service.Service, error) {
	svr := service.Service{}
	cli := application.Client{}
	if err := s.stmts[FindOneByID].QueryRow(sid).Scan(&svr.ID, &svr.Name, &svr.Description, &svr.Enabled, &svr.Status, &svr.StatusUpdateAt, &svr.Version, &svr.CreateAt, &cli.ClientID, &cli.ClientSecret); err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("service %s not find", sid)
		}

		return nil, exception.NewInternalServerError("query single service client error, %s", err)
	}

	svr.Client = &cli

	return &svr, nil
}

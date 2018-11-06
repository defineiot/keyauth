package mysql

import (
	"database/sql"
	"time"

	"github.com/satori/go.uuid"

	"github.com/defineiot/keyauth/dao/application"
	"github.com/defineiot/keyauth/internal/exception"
)

func (s *store) Registration(userID, name, description, website, clientID string) (*application.Application, error) {

	if userID == "" {
		return nil, exception.NewBadRequest("application user_id is missed")
	}
	if name == "" {
		return nil, exception.NewBadRequest("application name is missed")
	}

	ok, err := s.CheckAPPIsExistByName(userID, name)
	if err != nil {
		return nil, err
	}
	if ok {
		return nil, exception.NewBadRequest("application %s exist", name)
	}

	app := application.Application{ID: uuid.NewV4().String(), Name: name, UserID: userID, Website: website, Description: description, CreateAt: time.Now().Unix(), ClientID: clientID}

	_, err = s.stmts[CreateAPP].Exec(app.ID, app.Name, app.UserID, app.Website, app.LogoImage, app.Description, app.CreateAt, app.ClientID)
	if err != nil {
		return nil, exception.NewInternalServerError("insert application exec sql err, %s", err)
	}

	return &app, nil
}

func (s *store) CheckAPPIsExistByID(appID string) (bool, error) {
	var id string
	if err := s.stmts[CheckExistByID].QueryRow(appID).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, exception.NewInternalServerError("query application exist by id error, %s", err)
	}

	return true, nil
}

func (s *store) CheckAPPIsExistByName(userID, name string) (bool, error) {
	var id string
	if err := s.stmts[CheckExistByName].QueryRow(name, userID).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, exception.NewInternalServerError("query application exist by name error, %s", err)
	}

	return true, nil
}

func (s *store) Unregistration(id string) error {
	ok, err := s.CheckAPPIsExistByID(id)
	if err != nil {
		return err
	}
	if !ok {
		return exception.NewBadRequest("application %s not exist", id)
	}

	ret, err := s.stmts[DeleteAPP].Exec(id)
	if err != nil {
		return exception.NewInternalServerError("delete application exec sql error, %s", err)
	}
	count, err := ret.RowsAffected()
	if err != nil {
		return exception.NewInternalServerError("get delete row affected error, %s", err)
	}
	if count == 0 {
		return exception.NewBadRequest("application %s not exist", id)
	}

	return nil
}

func (s *store) ListApplications(userID string) ([]*application.Application, error) {
	rows, err := s.stmts[ListUserAPPS].Query(userID)
	if err != nil {
		return nil, exception.NewInternalServerError("query domain list error, %s", err)
	}
	defer rows.Close()

	applications := []*application.Application{}
	for rows.Next() {
		app := application.Application{}
		if err := rows.Scan(&app.ID, &app.Name, &app.UserID, &app.Website, &app.LogoImage, &app.Description, &app.CreateAt, &app.ClientID); err != nil {
			return nil, exception.NewInternalServerError("scan application record error, %s", err)
		}
		applications = append(applications, &app)
	}
	return applications, nil
}

func (s *store) GetApplication(appid string) (*application.Application, error) {
	app := application.Application{}
	err := s.stmts[GetUserAPP].QueryRow(appid).Scan(
		&app.ID, &app.Name, &app.UserID, &app.Website, &app.LogoImage, &app.Description, &app.CreateAt, &app.ClientID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("application %s not find", appid)
		}

		return nil, exception.NewInternalServerError("query single application error, %s", err)
	}

	return &app, nil
}

package mysql

import (
	"database/sql"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/defineiot/keyauth/dao/application"
	"github.com/defineiot/keyauth/internal/exception"
	"github.com/defineiot/keyauth/internal/tools"
)

func (s *store) CreateApplication(app *application.Application) error {
	if err := app.Validate(); err != nil {
		return err
	}

	ok, err := s.CheckAPPIsExistByName(app.UserID, app.Name)
	if err != nil {
		return err
	}
	if ok {
		return exception.NewBadRequest("application %s exist", app.Name)
	}

	app.ClientID = tools.MakeUUID(24)
	app.ClientSecret = tools.MakeUUID(32)
	app.CreateAt = time.Now().Unix()
	app.ID = uuid.NewV4().String()

	_, err = s.stmts[CreateAPP].Exec(app.ID, app.Name, app.UserID, app.Website, app.LogoImage,
		app.Description, app.CreateAt, app.RedirectURI, app.ClientID, app.ClientSecret,
		app.Locked, app.TokenExpireTime)
	if err != nil {
		return exception.NewInternalServerError("insert application exec sql err, %s", err)
	}

	return nil
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

func (s *store) DeleteApplication(appID string) error {
	ok, err := s.CheckAPPIsExistByID(appID)
	if err != nil {
		return err
	}
	if !ok {
		return exception.NewBadRequest("application %s not exist", appID)
	}

	ret, err := s.stmts[DeleteAPP].Exec(appID)
	if err != nil {
		return exception.NewInternalServerError("delete application exec sql error, %s", err)
	}
	count, err := ret.RowsAffected()
	if err != nil {
		return exception.NewInternalServerError("get delete row affected error, %s", err)
	}
	if count == 0 {
		return exception.NewBadRequest("application %s not exist", appID)
	}

	return nil
}

func (s *store) ListUserApplications(userID string) ([]*application.Application, error) {
	rows, err := s.stmts[ListUserAPPS].Query(userID)
	if err != nil {
		return nil, exception.NewInternalServerError("query domain list error, %s", err)
	}
	defer rows.Close()

	applications := []*application.Application{}
	for rows.Next() {
		app := new(application.Application)
		if err := rows.Scan(&app.ID, &app.Name, &app.UserID, &app.Website, &app.LogoImage,
			&app.Description, &app.CreateAt, &app.RedirectURI, &app.ClientID, &app.ClientSecret,
			&app.Locked, &app.LastLoginTime, &app.LastLoginIP, &app.LoginFailedTimes,
			&app.LoginSuccessTimes, &app.TokenExpireTime); err != nil {
			return nil, exception.NewInternalServerError("scan application record error, %s", err)
		}
		applications = append(applications, app)
	}
	return applications, nil
}

func (s *store) GetApplication(appid string) (*application.Application, error) {
	app := application.Application{}
	err := s.stmts[GetUserAPP].QueryRow(appid).Scan(&app.ID, &app.Name, &app.UserID, &app.Website,
		&app.LogoImage, &app.Description, &app.CreateAt, &app.RedirectURI, &app.ClientID,
		&app.ClientSecret, &app.Locked, &app.LastLoginTime, &app.LastLoginIP, &app.LoginFailedTimes,
		&app.LoginSuccessTimes, &app.TokenExpireTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("application %s not find", appid)
		}

		return nil, exception.NewInternalServerError("query single application error, %s", err)
	}

	return &app, nil
}

func (s *store) GetApplicationByClientID(clientID string) (*application.Application, error) {
	app := application.Application{}
	err := s.stmts[GetUserAPPByClientID].QueryRow(clientID).Scan(&app.ID, &app.Name, &app.UserID,
		&app.Website, &app.LogoImage, &app.Description, &app.CreateAt, &app.RedirectURI,
		&app.ClientID, &app.ClientSecret, &app.Locked, &app.LastLoginTime, &app.LastLoginIP,
		&app.LoginFailedTimes, &app.LoginSuccessTimes, &app.TokenExpireTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("application %s not find", clientID)
		}

		return nil, exception.NewInternalServerError("query single application error, %s", err)
	}

	return &app, nil
}

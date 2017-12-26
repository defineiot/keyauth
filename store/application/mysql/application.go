package mysql

import (
	"database/sql"
	"math/rand"
	"strings"
	"time"

	"github.com/satori/go.uuid"

	"openauth/api/exception"
	"openauth/store/application"
)

func (s *store) Registration(userID, name, redirectURI, clientType, description, website string) (*application.Application, error) {

	if userID == "" {
		return nil, exception.NewBadRequest("application user_id is missed")
	}
	if name == "" {
		return nil, exception.NewBadRequest("application name is missed")
	}
	if clientType != "confidential" && clientType != "public" {
		return nil, exception.NewBadRequest("application's client_type must one of confidential or public")
	}

	ok, err := s.CheckAPPIsExistByName(userID, name)
	if err != nil {
		return nil, err
	}
	if ok {
		return nil, exception.NewBadRequest("application %s exist", name)
	}

	clientID, err := makeuuid(24)
	if err != nil {
		return nil, exception.NewInternalServerError("initial client_id error, %s", err)
	}
	clientSecret, err := makeuuid(32)
	if err != nil {
		return nil, exception.NewInternalServerError("initial client_secret error, %s", err)
	}

	client := application.Client{ClientID: clientID, ClientSecret: clientSecret, ClientType: clientType, RedirectURI: redirectURI}
	app := application.Application{ID: uuid.NewV4().String(), Name: name, UserID: userID, Website: website, Description: description, CreateAt: time.Now().Unix()}
	app.Client = &client

	_, err = s.stmts[CreateAPP].Exec(app.ID, app.Name, app.UserID, client.ClientID, client.ClientSecret, client.ClientType, app.Website, app.LogoImage, app.Description, client.RedirectURI, app.CreateAt)
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

func (s *store) GetUserApps(userID string) ([]*application.Application, error) {
	rows, err := s.stmts[GetUserAPPS].Query(userID)
	if err != nil {
		return nil, exception.NewInternalServerError("query domain list error, %s", err)
	}
	defer rows.Close()

	applications := []*application.Application{}
	for rows.Next() {
		app := application.Application{}
		cli := application.Client{}
		if err := rows.Scan(&app.ID, &app.Name, &app.UserID, &cli.ClientID, &cli.ClientSecret, &cli.ClientType, &app.Website, &app.LogoImage, &app.Description, &cli.RedirectURI, &app.CreateAt); err != nil {
			return nil, exception.NewInternalServerError("scan application record error, %s", err)
		}
		app.Client = &cli
		applications = append(applications, &app)
	}
	return applications, nil
}

func (s *store) GetClient(clientID string) (*application.Client, error) {
	cli := new(application.Client)
	err := s.stmts[GetClient].QueryRow(clientID).Scan(&cli.ClientID, &cli.ClientSecret, &cli.ClientType, &cli.RedirectURI)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("client %s not find", clientID)
		}

		return nil, exception.NewInternalServerError("query single application's client error, %s", err)
	}

	return cli, nil
}

func makeuuid(lenth int) (string, error) {
	charlist := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	password := make([]string, lenth)
	rand.Seed(time.Now().UnixNano() + int64(lenth))
	for i := 0; i < lenth; i++ {
		rn := rand.Intn(len(charlist))
		w := charlist[rn : rn+1]
		password = append(password, w)
	}

	return strings.Join(password, ""), nil
}

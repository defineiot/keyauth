package mysql

import (
	"database/sql"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/satori/go.uuid"

	"openauth/api/exception"
	"openauth/pkg/application"
)

var (
	createStmt *sql.Stmt
	deleteStmt *sql.Stmt
)

// NewApplicationService use to new application instance
func NewApplicationService(db *sql.DB) application.Service {
	return &manager{db: db}
}

type manager struct {
	db *sql.DB
}

func (m *manager) Registration(userID, name, redirectURI, clientType, description, website string) (*application.Application, error) {
	var (
		once sync.Once
		err  error
	)

	if userID == "" {
		return nil, exception.NewBadRequest("application user_id is missed")
	}
	if name == "" {
		return nil, exception.NewBadRequest("application name is missed")
	}
	if clientType != "confidential" && clientType != "public" {
		return nil, exception.NewBadRequest("application's client_type must one of confidential or public")
	}

	ok, err := m.CheckAPPIsExistByName(userID, name)
	if err != nil {
		return nil, err
	}
	if ok {
		return nil, exception.NewBadRequest("application %s exist", name)
	}

	once.Do(func() {
		createStmt, err = m.db.Prepare("INSERT INTO `application` (id, name, user_id, client_id, client_secret, client_type, website, logo_image, description, redirect_uri, create_at) VALUES (?,?,?,?,?,?,?,?,?,?,?)")
	})
	if err != nil {
		return nil, exception.NewInternalServerError("prepare insert application stmt error, domain: %s, %s", name, err)
	}

	clientID, err := makeuuid(10)
	if err != nil {
		return nil, exception.NewInternalServerError("initial client_id error, %s", err)
	}
	clientSecret, err := makeuuid(16)
	if err != nil {
		return nil, exception.NewInternalServerError("initial client_secret error, %s", err)
	}

	app := application.Application{ID: uuid.NewV4().String(), Name: name, UserID: userID, ClientID: clientID, ClientSecret: clientSecret, ClientType: clientType, Website: website, Description: description, RedirectURI: redirectURI, CreateAt: time.Now().Unix()}
	_, err = createStmt.Exec(app.ID, app.Name, app.UserID, app.ClientID, app.ClientSecret, app.ClientType, app.Website, app.LogoImage, app.Description, app.RedirectURI, app.CreateAt)
	if err != nil {
		return nil, exception.NewInternalServerError("insert application exec sql err, %s", err)
	}

	return &app, nil
}

func (m *manager) CheckAPPIsExistByID(appID string) (bool, error) {
	var id string
	if err := m.db.QueryRow("SELECT id FROM application WHERE id = ?", appID).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, exception.NewInternalServerError("query application exist by id error, %s", err)
	}

	return true, nil
}

func (m *manager) CheckAPPIsExistByName(userID, name string) (bool, error) {
	var id string
	if err := m.db.QueryRow("SELECT id FROM application WHERE name = ? AND user_id = ?", name, userID).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, exception.NewInternalServerError("query application exist by name error, %s", err)
	}

	return true, nil
}

func (m *manager) Unregistration(id string) error {
	var (
		once sync.Once
		err  error
	)

	ok, err := m.CheckAPPIsExistByID(id)
	if err != nil {
		return err
	}
	if !ok {
		return exception.NewBadRequest("application %s not exist", id)
	}

	once.Do(func() {
		deleteStmt, err = m.db.Prepare("DELETE FROM application WHERE id = ?")
	})
	if err != nil {
		return exception.NewInternalServerError("prepare delete application stmt error, %s", err)
	}

	ret, err := deleteStmt.Exec(id)
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

func (m *manager) GetUserApps(userID string) ([]*application.Application, error) {
	rows, err := m.db.Query("SELECT id,name,user_id,client_id,client_secret,client_type,website,logo_image,description,redirect_uri,create_at FROM application WHERE user_id = ? ORDER BY create_at DESC", userID)
	if err != nil {
		return nil, exception.NewInternalServerError("query domain list error, %s", err)
	}
	defer rows.Close()

	applications := []*application.Application{}
	for rows.Next() {
		app := application.Application{}
		if err := rows.Scan(&app.ID, &app.Name, &app.UserID, &app.ClientID, &app.ClientSecret, &app.ClientType, &app.Website, &app.LogoImage, &app.Description, &app.RedirectURI, &app.CreateAt); err != nil {
			return nil, exception.NewInternalServerError("scan application record error, %s", err)
		}
		applications = append(applications, &app)
	}
	return applications, nil
}

func makeuuid(lenth int) (string, error) {
	charlist := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	password := make([]string, lenth)
	rand.Seed(time.Now().Unix() + int64(lenth))
	for i := 0; i < lenth; i++ {
		rn := rand.Intn(len(charlist))
		w := charlist[rn : rn+1]
		password = append(password, w)
	}

	return strings.Join(password, ""), nil
}

package mysql

import (
	"database/sql"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/satori/go.uuid"

	"openauth/api/exception"
	"openauth/storage/application"
)

var (
	createStmt *sql.Stmt
	deleteStmt *sql.Stmt
)

// NewApplicationStorage use to new application instance
func NewApplicationStorage(db *sql.DB) application.Storage {
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

	return nil, nil
}

func (m *manager) Unregistration(userID, clientID string) error {
	return nil
}

func (m *manager) GetUserApps(userID string) ([]*application.Application, error) {
	return nil, nil
}

func makeuuid(lenth int) (string, error) {
	charlist := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	password := make([]string, lenth)
	rand.Seed(time.Now().Unix())
	for i := 0; i < lenth; i++ {
		rn := rand.Intn(lenth)
		w := charlist[rn : rn+1]
		password = append(password, w)
	}

	return strings.Join(password, ""), nil
}

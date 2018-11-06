package application

import (
	"github.com/defineiot/keyauth/dao/client"
)

// Application is oauth2's client: https://tools.ietf.org/html/rfc6749#section-2
type Application struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Website     string `json:"website,omitempty"`
	LogoImage   string `json:"logo_image,omitempty"`
	Description string `json:"description"`
	CreateAt    int64  `json:"create_at"`
	ClientID    string `json:"-"`
	*client.Client
}

// Store application storage
type Store interface {
	Reader
	Writer
	Close() error
}

// Reader use to read application information from store
type Reader interface {
	ListApplications(userID string) ([]*Application, error)
	GetApplication(appid string) (*Application, error)

	CheckAPPIsExistByID(appID string) (bool, error)
	CheckAPPIsExistByName(userID, name string) (bool, error)
}

// Writer use to write application information from store
type Writer interface {
	Registration(userID, name, description, website, clientID string) (*Application, error)
	Unregistration(id string) error
}

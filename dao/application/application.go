package application

import (
	"github.com/defineiot/keyauth/dao/models"
)

// Store application storage
type Store interface {
	Reader
	Writer
	Close() error
}

// Reader use to read application information from store
type Reader interface {
	ListUserApplications(userID string) ([]*models.Application, error)
	GetApplication(appID string) (*models.Application, error)
	GetApplicationByClientID(clientID string) (*models.Application, error)

	CheckAPPIsExistByID(appID string) (bool, error)
	CheckAPPIsExistByName(userID, name string) (bool, error)
}

// Writer use to write application information from store
type Writer interface {
	CreateApplication(app *models.Application) error
	DeleteApplication(appID string) error
}

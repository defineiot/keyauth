package service

import (
	"openauth/store/application"
)

// Service is service provider
type Service struct {
	ID             string              `json:"id"`
	Name           string              `json:"name"`
	Description    string              `json:"description"`
	Version        string              `json:"version"`
	Enabled        bool                `json:"enabled"`
	Status         string              `json:"status"`
	StatusUpdateAt int64               `json:"status_update_at"`
	CreateAt       int64               `json:"create_at"`
	Client         *application.Client `json:"client"`
}

// Store service store interface
type Store interface {
	StoreReader
	StoreWriter
	Close() error
}

// StoreReader read service information from store
type StoreReader interface {
	FindAllService() ([]*Service, error)
	FindServiceByID() (*Service, error)
}

// StoreWriter write service information to store
type StoreWriter interface {
	SaveService(name, description string) (*Service, error)
	DeleteService(sid string) error
}

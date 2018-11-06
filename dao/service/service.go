package service

import (
	"github.com/defineiot/keyauth/dao/client"
)

// Service is service provider
type Service struct {
	Name           string     `json:"name"`
	Description    string     `json:"description,omitempty"`
	Version        string     `json:"version,omitempty"`
	Enabled        bool       `json:"enabled"`
	Status         string     `json:"status,omitempty"`
	StatusUpdateAt int64      `json:"status_update_at,omitempty"`
	CreateAt       int64      `json:"create_at"`
	ClientID       string     `json:"-"`
	Features       []*Feature `json:"features,omitempty"`
	*client.Client
}

// Feature Service's features
type Feature struct {
	ID                 int64  `json:"id"`
	Name               string `json:"name"`
	Method             string `json:"method"`
	Endpoint           string `json:"endpoint"`
	Description        string `json:"description,omitempty"`
	IsDeleted          bool   `json:"is_deleted,omitempty"`
	WhenDeletedVersion string `json:"when_deleted_version,omitempty"`
	IsAdded            bool   `json:"is_added,omitempty"`
	WhenAddedVersion   string `json:"when_added_version,omitempty"`
	ServiceName        string `json:"service_name"`
}

// Store service store interface
type Store interface {
	Reader
	Writer
	Close() error
}

// Reader read service information from store
type Reader interface {
	ListServices() ([]*Service, error)
	GetService(name string) (*Service, error)
	CheckServiceIsExist(name string) (bool, error)
	ListServiceFeatures(name string) ([]*Feature, error)
	ListDomainFeatures() ([]*Feature, error)
	ListMemberFeatures() ([]*Feature, error)
	ListRoleFeatures(name string) ([]*Feature, error)
	CheckFeatureIsExist(featureID int64) (bool, error)
	CheckServiceHasFeature(serviceName, featureName string) (bool, error)
	GetServiceByClientID(clientID string) (*Service, error)
}

// Writer write service information to store
type Writer interface {
	CreateService(name, description, clientID string) (*Service, error)
	RegistryServiceFeatures(name string, features ...Feature) error

	DeleteService(name string) error
}

package service

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/defineiot/keyauth/internal/exception"
)

const (
	HTTP_API_SERVICE Type = iota
	AGENT_WORKER_SERVICE
)

// Type 服务类型
type Type int

// Feature Service's features
type Feature struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Method             string `json:"method"`
	Endpoint           string `json:"endpoint"`
	Description        string `json:"description,omitempty"`
	IsDeleted          bool   `json:"is_deleted,omitempty"`
	WhenDeletedVersion string `json:"when_deleted_version,omitempty"`
	IsAdded            bool   `json:"is_added,omitempty"`
	WhenAddedVersion   string `json:"when_added_version,omitempty"`
	ServiceID          string `json:"service_id,omitempty"`
}

// Service is service provider
type Service struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	Description    string     `json:"description,omitempty"`
	Enabled        bool       `json:"enabled"`
	Status         string     `json:"status,omitempty"`
	StatusUpdateAt int64      `json:"status_update_at,omitempty"`
	Version        string     `json:"version,omitempty"`
	CreateAt       int64      `json:"create_at"`
	ClientID       string     `json:"client_id"`
	ClientSecret   string     `json:"client_secret"`
	Type           Type       `json:"service_type"`
	Features       []*Feature `json:"features,omitempty"`
}

func (s *Service) String() string {
	str, err := json.Marshal(s)
	if err != nil {
		log.Printf("E! marshal role to string error: %s", err)
		return fmt.Sprintf("ID: %s, Name: %s", s.ID, s.Name)
	}

	return string(str)
}

// Validate 服务创建检查
func (s *Service) Validate() error {
	if s.Name == "" {
		exception.NewBadRequest("the service's name is required!")
	}

	return nil
}

// Store service store interface
type Store interface {
	Reader
	Writer
	Close() error
}

// Reader read service information from store
type Reader interface {
	CheckInstanceExist(instanceName, serviceName string) (bool, error)

	ListServices() ([]*Service, error)
	GetService(name string) (*Service, error)
	CheckServiceIsExist(name string) (bool, error)
	ListServiceFeatures(name string) ([]*Feature, error)

	ListDomainFeatures() ([]*Feature, error)
	ListMemberFeatures() ([]*Feature, error)
	ListRoleFeatures(name string) ([]*Feature, error)
	CheckFeatureIsExist(featureID int64) (bool, error)
	CheckServiceHasFeature(serviceName, featureName string) (bool, error)
}

// Writer write service information to store
type Writer interface {
	CreateService(service *Service) error
	RegistryServiceFeatures(name string, features ...Feature) error
	DeleteService(id string) error

	UpdateInstanceOffline(instanceName, serviceName string) error
}

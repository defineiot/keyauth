package service

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/defineiot/keyauth/internal/exception"
)

// Instance 服务实例
type Instance struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Address     string `json:"address"`
	Version     string `json:"version"`
	GitBranch   string `json:"git_branch"`
	GitCommit   string `json:"git_commit"`
	BuildEnv    string `json:"build_env"`
	BuildAt     string `json:"build_at"`
	ServiceName string `json:"service_name"`
	ServiceType string `json:"service_type"`

	Status  string `json:"status"`
	Online  int64  `json:"online_at"`
	Offline int64  `json:"offline_at"`
}

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
	ServiceID          string `json:"service_id"`
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
	FindAllInstances() ([]*Instance, error)

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

	SaveInstance(*Instance) error
	UpdateInstance(instance *Instance) error
	UpdateInstanceOffline(instanceName, serviceName string) error
}

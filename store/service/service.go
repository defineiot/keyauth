package service

// Service is service provider
type Service struct {
	Features   []string
	Status     string
	Components []Component
	Active     bool
}

// Component is service sub service
type Component struct {
	Name      string   `json:"name"`
	Instances []string `json:"instances"`
	Error     string   `json:"error"`
}

// Store service store interface
type Store interface {
	SaveService() (*Service, error)
	UpdateService() (*Service, error)
	DeleteService() (*Service, error)
	FindAllService() ([]*Service, error)
	FindServiceByID() (*Service, error)
}

// Manager is catalog service
type Manager interface {
	// List all Servers and the Description
	ListServer() ([]string, error)
	// List all feature, return function list
	ListFeature() ([]string, error)
	// Verify that this feature is available
	HaveFeatures(features ...string) (bool, error)
	// Service Catalog
	Catalog() (*[]map[string]Service, error)
}

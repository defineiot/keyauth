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

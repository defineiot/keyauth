package register

import (
	"errors"
	"time"
)

// ServiceInstance todo
type ServiceInstance struct {
	Name        string `json:"name"`
	ServiceName string `json:"service_name"`
	Type        string `json:"type"`
	Address     string `json:"address"`
	Version     string `json:"version"`
	GitBranch   string `json:"git_branch"`
	GitCommit   string `json:"git_commit"`
	BuildEnv    string `json:"build_env"`
	BuildAt     string `json:"build_at"`

	Prefix   string        `json:"-"`
	Interval time.Duration `json:"-"`
	TTL      int           `json:"-"`
}

// Validate 服务实例注册参数校验
func (s *ServiceInstance) Validate() error {
	if s.Name == "" && s.ServiceName == "" || s.Type == "" {
		return errors.New("service instance name or service_name or type not config")
	}

	return nil
}

// Register 服务注册接口
type Register interface {
	Registe(service *ServiceInstance) error
	UnRegiste() error
}

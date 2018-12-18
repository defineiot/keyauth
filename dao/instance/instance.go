package instance

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/defineiot/keyauth/internal/exception"

	"github.com/defineiot/keyauth/dao/service"
)

// Instance 服务实例
type Instance struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	Version   string `json:"version"`
	GitBranch string `json:"git_branch"`
	GitCommit string `json:"git_commit"`
	BuildEnv  string `json:"build_env"`
	BuildAt   string `json:"build_at"`
	ServiceID string `json:"service_id,omitempty"`
	OnlineAt  int64  `json:"online_at"`
	OfflineAt int64  `json:"offline_at"`

	Features []*service.Feature `json:"features"`
}

func (i *Instance) String() string {
	str, err := json.Marshal(i)
	if err != nil {
		log.Printf("E! marshal domain to string error: %s", err)
		return fmt.Sprintf("ID: %s, Name: %s", i.ID, i.Name)
	}

	return string(str)
}

// Validate 注册校验
func (i *Instance) Validate() error {
	if i.ServiceID == "" {
		return exception.NewBadRequest("the instance service id is required!")
	}

	if i.Name == "" {
		return exception.NewBadRequest("the instance name is required!")
	}

	if i.Address == "" {
		return exception.NewBadRequest("the instance address is required!")
	}

	return nil
}

// Store is project service
type Store interface {
	Reader
	Writer
	Close() error
}

// Reader read project information from store
type Reader interface {
	ListServiceInstances(serviceID string) ([]*Instance, error)
}

// Writer for write data to store
type Writer interface {
	RegistryServiceInstance(ins *Instance) error
	UnRegisteServiceInstance(serviceID, instanceID string) error
}

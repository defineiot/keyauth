package instance

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/defineiot/keyauth/dao/models"
	"github.com/defineiot/keyauth/internal/exception"
)

// Instance 服务实例
type Instance struct {
	ID        string `json:"id,omitempty"`         // 实例的ID, 比如: web-aabbccxxxdd(项目代号-8位随机字符串)
	Name      string `json:"name"`                 // 实例名称
	Address   string `json:"address"`              // 实例运行的地址, 比如:   127.0.0.1:[8080/8081/8082], 如果不需要监听端口端口列表可以省略
	Version   string `json:"version"`              // 服务的版本号, 采用GUN风格版本命名风格: 主版本号.子版本号[.修正版本号[.编译版本号]] 标准格式为: 1.1.12.58
	GitBranch string `json:"git_branch"`           // git对于的分支信息
	GitCommit string `json:"git_commit"`           // 具体的commit id
	BuildEnv  string `json:"build_env"`            // 编辑工具的信息, 比如 Go 1.10.4/ JDK_8
	BuildAt   string `json:"build_at"`             // 编译打包的时间
	ServiceID string `json:"service_id,omitempty"` // 实例所属的服务ID
	OnlineAt  int64  `json:"online_at"`            // 实例注册上线的时间
	OfflineAt int64  `json:"offline_at"`           // 实例下线的时间

	Features []*models.Feature `json:"features"` // 实例提供的功能列表
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

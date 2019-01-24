package service

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/defineiot/keyauth/internal/exception"
)

const (
	Internal Type = "internal_rpc"     // 内部调用的控制面类型的服务, 提供了RPC能力,需要注册到API 网关对内提供服务
	Public        = "controlle_pannel" // 需要对外发布的控制面类型的服务, 提供了RPC能力, 需要注册到API 网关对外提供服务
	Agent         = "data_pannel"      // 数据面类型的服务, 不提供RPC能力

	Unknown     Status = "unknown"     // 刚创建好服务, 并没有服务实例启动
	Unavailable        = "unavailable" // 后端服务没有实例提供服务
	Avaliable          = "avaliable"   // 后台服务有实例提供服务
	Upgrading          = "upgrading"   // 多个服务实例版本不一致, 处于升级状态
	Downgrading        = "downgrading" // 多个服务实例版本不一致, 处于回滚状态
)

// Type 服务类型
type Type string

// Status 服务状态
type Status string

// Feature Service's features
type Feature struct {
	ID             string `json:"id"`                             // 功能唯一标识符
	Name           string `json:"name"`                           // 功能的名称
	Tag            string `json:"tag,omitempty"`                  // 功能的标签, 如果该功能对应的HTTP类型RPC, 标签可以为 POST/GET/DELETE
	HTTPEndpoint   string `json:"endpoint,omitempty"`             // 该功能对应的HTTP类型RPC, 比如 /<service_name>/<resource_name>/<action>
	Description    string `json:"description,omitempty"`          // 该功能的描述信息
	IsDeleted      bool   `json:"is_deleted,omitempty"`           // 该功能是否已经废弃
	DeletedVersion string `json:"when_deleted_version,omitempty"` // 该功能在那个版本废弃的
	DeleteAt       int64  `json:"when_deleted_time"`              // 功能废弃的时间
	IsAdded        bool   `json:"is_added,omitempty"`             // 该功能是否是新增功能
	AddedVersion   string `json:"when_added_version,omitempty"`   // 该功能在那个版本新增的
	AddedAt        int64  `json:"when_added_time"`                // 功能注册的时间
	ServiceID      string `json:"service_id,omitempty"`           // 该功能属于那个服务
}

// Service is service provider
type Service struct {
	ID               string     `json:"id"`                          // 唯一ID
	Type             Type       `json:"type"`                        // 服务类型
	Name             string     `json:"name"`                        // 名称
	Description      string     `json:"description,omitempty"`       // 描述信息
	Enabled          bool       `json:"enabled"`                     // 是否启用该服务
	Status           Status     `json:"status,omitempty"`            // 服务状态(unavailable/avaliable/upgrading/downgrading)
	StatusUpdateAt   int64      `json:"status_update_at,omitempty"`  // 状态更新的时间
	CurrentVersion   string     `json:"current_version,omitempty"`   // 当前版本信息, 通过比较获取的实例对应的版本
	UpgradeVersion   string     `json:"upgrade_version,omitempty"`   // 如果服务实例版本不一致时, 新注册的实例的版本大于当前版本, 则记录该升级版本的信息
	DowngradeVersion string     `json:"downgrade_version,omitempty"` // 如果服务实例版本不一致时, 新注册的实例的版本小于当前版本, 则记录该升级版本的信息
	CreateAt         int64      `json:"create_at"`                   // 创建的时间
	UpdateAt         int64      `json:"update_at"`                   // 更新时间
	ClientID         string     `json:"client_id"`                   // 客户端ID
	ClientSecret     string     `json:"client_secret"`               // 客户端秘钥
	TokenExpireTime  int64      `json:"token_expire_time"`           // 客户端凭证申请的token的过期时间
	Features         []*Feature `json:"features,omitempty"`          // 服务功能列表
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
		return exception.NewBadRequest("the service's name is required!")
	}

	if s.Type == "" {
		return exception.NewBadRequest("the service's type is required!")
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
	ListServices() ([]*Service, error)
	GetServiceByID(id string) (*Service, error)
	GetServiceByName(name string) (*Service, error)
	GetServiceByClientID(clientID string) (*Service, error)

	ListServiceFeatures(serviceID string) ([]*Feature, error)
	ListRoleFeatures(roleID string) ([]*Feature, error)
}

// Writer write service information to store
type Writer interface {
	CreateService(service *Service) error
	DeleteService(id string) error

	RegistryServiceFeatures(serviceID, version string, features ...*Feature) error
	AssociateFeaturesToRole(roleID string, features ...*Feature) error
	UnlinkFeatureFromRole(roleID string, features ...*Feature) error
}

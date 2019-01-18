package role

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/defineiot/keyauth/dao/service"
	"github.com/defineiot/keyauth/internal/exception"
)

// Role is rbac's role
type Role struct {
	ID          string `json:"id"`
	Name        string `json:"name"`        // 角色名称
	Description string `json:"description"` // 角色描述
	CreateAt    int64  `json:"create_at"`   // 创建时间
	UpdateAt    int64  `json:"update_at"`   // 更新时间

	Features []*service.Feature `json:"features"` // 角色的功能列表
}

func (r *Role) String() string {
	str, err := json.Marshal(r)
	if err != nil {
		log.Printf("E! marshal role to string error: %s", err)
		return fmt.Sprintf("ID: %s, Name: %s", r.ID, r.Name)
	}

	return string(str)
}

// Validate 校验检查
func (r *Role) Validate() error {
	if r.Name == "" {
		return exception.NewBadRequest("the role's name is required!")
	}

	return nil
}

// Store is an role service
type Store interface {
	Reader
	Writer
	Close() error
}

// Reader for read data from store
type Reader interface {
	GetRole(name string) (*Role, error)
	CheckRoleExist(name string) (bool, error)
	ListRole() ([]*Role, error)
	// VerifyRole(name string, feature string) (bool, error)
}

// Writer for write data to store
type Writer interface {
	CreateRole(role *Role) error
	UpdateRole(name, description string) (*Role, error)
	DeleteRole(name string) error
}

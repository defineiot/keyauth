package models

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/defineiot/keyauth/internal/exception"
)

// Role is rbac's role
type Role struct {
	ID          string `json:"id"`
	Name        string `json:"name"`                  // 角色名称
	Description string `json:"description,omitempty"` // 角色描述
	CreateAt    int64  `json:"create_at,omitempty"`   // 创建时间
	UpdateAt    int64  `json:"update_at,omitempty"`   // 更新时间

	Features []*Feature `json:"features,omitempty"` // 角色的功能列表
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

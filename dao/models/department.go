package models

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/defineiot/keyauth/internal/exception"
)

// Department user's department
type Department struct {
	ID       string `json:"id"`
	Number   string `json:"number,omitempty"`    // 部门编号
	Name     string `json:"name,omitempty"`      // 部门名称
	Grade    string `json:"grade,omitempty"`     // 第几级部门
	Path     string `json:"path,omitempty"`      // 部门访问路径
	CreateAt int64  `json:"create_at,omitempty"` // 部门创建时间
	DomainID string `json:"domain_id,omitempty"` // 部门所属域

	ParentID  string `json:"parent_id,omitempty"`  // 上级部门ID
	ManagerID string `json:"manager_id,omitempty"` // 部门管理者ID

	Users    []*User    `json:"users,omitempty"`    // 部门用户
	Projects []*Project `json:"projects,omitempty"` // 部门可以访问的项目
	Roles    []*Role    `json:"roles,omitempty"`    // 部门人员的角色
}

func (d *Department) String() string {
	str, err := json.Marshal(d)
	if err != nil {
		log.Printf("E! marshal department to string error: %s", err)
		return fmt.Sprintf("ID: %s, Name: %s", d.ID, d.Name)
	}

	return string(str)
}

// Validate 校验创建的数据
func (d *Department) Validate() error {
	if d.Name == "" {
		return exception.NewBadRequest("the department's name is required!")
	}

	if d.DomainID == "" {
		return exception.NewBadRequest("the department's domain id is required!")
	}

	if len(d.Name) > 128 {
		return exception.NewBadRequest("department's name is too long,  max length is 128")
	}

	return nil
}

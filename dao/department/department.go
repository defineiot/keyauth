package department

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/defineiot/keyauth/internal/exception"
)

// Department user's department
type Department struct {
	ID       string `json:"id"`
	Number   string `json:"-"`         // 部门编号
	Name     string `json:"name"`      // 部门名称
	Grade    string `json:"grade"`     // 第几级部门
	Path     string `json:"path"`      // 部门访问路径
	CreateAt int64  `json:"create_at"` // 部门创建时间
	DomainID string `json:"domain_id"` // 部门所属域

	ParentID  string `json:"parent_id"`  // 上级部门ID
	ManagerID string `json:"manager_id"` // 部门管理者ID
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

// Store is user service
type Store interface {
	Reader
	Writer
	Close() error
}

// Reader use to read department information form store
type Reader interface {
	GetDepartment(depID string) (*Department, error)
	ListSubDepartments(parentDepID string) ([]*Department, error)
}

// Writer use to write department information to store
type Writer interface {
	CreateDepartment(d *Department) (*Department, error)
	DelDepartment(depID string) error
}

package models

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/defineiot/keyauth/internal/exception"
)

// Project tenant resource container
type Project struct {
	ID          string  `json:"id"`                    // 项目唯一ID
	Name        string  `json:"name"`                  // 项目名称
	Picture     string  `json:"picture,omitempty"`     // 项目描述图片
	Latitude    float64 `json:"latitude,omitempty"`    // 项目所处地理位置的经度
	Longitude   float64 `json:"longitude,omitempty"`   // 项目所处地理位置的纬度
	Enabled     bool    `json:"enabled,omitempty"`     // 禁用项目, 该项目所有人暂时都无法访问
	Owner       string  `json:"owner_id,omitempty"`    // 项目所有者, 一般为创建者
	Description string  `json:"description,omitempty"` // 项目描述
	DomainID    string  `json:"domain_id,omitempty"`   // 所属域ID
	CreateAt    int64   `json:"create_at,omitempty"`   // 创建时间
	UpdateAt    int64   `json:"update_at,omitempty"`   // 项目修改时间
}

func (p *Project) String() string {
	str, err := json.Marshal(p)
	if err != nil {
		log.Printf("E! marshal project to string error: %s", err)
		return fmt.Sprintf("ID: %s, Name: %s", p.ID, p.Name)
	}

	return string(str)
}

// Validate todo
func (p *Project) Validate() error {
	if p.Name == "" {
		return exception.NewBadRequest("project's name is required!")
	}

	if len(p.Name) > 128 {
		return exception.NewBadRequest("project's name is too long,  max length is 128")
	}

	if p.Owner == "" {
		return exception.NewBadRequest("project's owner_id is required!")
	}

	if p.DomainID == "" {
		return exception.NewBadRequest("project's domain_id required!")
	}

	return nil
}

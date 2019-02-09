package models

import (
	"fmt"

	"github.com/defineiot/keyauth/internal/exception"
)

const (
	RegistryCode CodePurpose = iota
	ChangePasswordCode
	LoginCode

	SendByMobile SendMode = iota // 通过手机发放验证码
	SendByEmail                  // 通过邮件发送验证码
	SendByAPP                    // 通过手机应用(app)发放验证码
)

// CodePurpose 验证码的用途
type CodePurpose int

// SendMode 验证码的发送方式
type SendMode int

// VerifyCode 一次性使用验证码, 使用过后则立即删除
type VerifyCode struct {
	Code       string      `json:"code"`                  // 验证码, 6位的数字
	Purpose    CodePurpose `json:"purpose,omitempty"`     // 用途
	SendMode   SendMode    `json:"send_mode,omitempty"`   // 发送方式
	SendTarget string      `json:"send_target,omitempty"` // 发送人地址
	CreateAt   int64       `json:"create_at,omitempty"`   // 创建时间
	ExpireAt   int64       `json:"expire_at,omitempty"`   // 过期时间
}

func (v *VerifyCode) String() string {
	return fmt.Sprintf("code: %s, purpose: %d, mode: %d, target: %s, create_at: %d, expire_at: %d",
		v.Code, v.Purpose, v.SendMode, v.SendTarget, v.CreateAt, v.ExpireAt)
}

// Validate 创建时校验
func (v *VerifyCode) Validate() error {
	if v.SendTarget == "" {
		return exception.NewBadRequest("send verify code, but the sending_target is empty!")
	}

	return nil
}

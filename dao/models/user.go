package models

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/defineiot/keyauth/internal/exception"
)

const (
	UserIDIndex  FoundUserIndex = "uid"
	AccountIndex FoundUserIndex = "account"
	MobileIndex  FoundUserIndex = "mobile"
	EmailIndex   FoundUserIndex = "email"
)

// FoundUserIndex 获取用户的索引类型
type FoundUserIndex string

// User info
type User struct {
	ID                string    `json:"id"`                            // 用户UUID
	Account           string    `json:"account,omitempty"`             // 用户账号名称
	Mobile            string    `json:"mobile,omitempty"`              // 手机号码, 用户可以通过手机进行注册和密码找回, 还可以通过手机号进行登录
	Email             string    `json:"email,omitempty"`               // 邮箱, 用户可以通过邮箱进行注册和照明密码
	Phone             string    `json:"phone,omitempty"`               // 用户的座机号码
	Address           string    `json:"address,omitempty"`             // 用户住址
	RealName          string    `json:"real_name,omitempty"`           // 用户真实姓名
	NickName          string    `json:"nick_name,omitempty"`           // 用户昵称, 用于在界面进行展示
	Gender            string    `json:"gender,omitempty"`              // 性别
	Avatar            string    `json:"avatar,omitempty"`              // 头像
	Language          string    `json:"language,omitempty"`            // 用户使用的语言
	City              string    `json:"city,omitempty"`                // 用户所在的城市
	Province          string    `json:"province,omitempty"`            // 用户所在的省
	Locked            string    `json:"locked,omitempty"`              // 是否冻结次用户
	CreateAt          int64     `json:"create_at,omitempty"`           // 用户创建的时间
	ExpiresActiveDays int       `json:"expires_active_days,omitempty"` // 用户多久未登录时(天), 冻结改用户, 防止僵尸用户的账号被利用
	Password          *Password `json:"password,omitempty"`            // 密码相关信息

	Domain         *Domain      `json:"domain,omitempty"`          // 如果需要对象由上层进行查找
	DefaultProject *Project     `json:"default_project,omitempty"` //  如果需要对象由上层进行查找
	Department     *Department  `json:"department,omitempty"`      // 所属部门信息
	LoginStatus    *LoginStatus `json:"login_status,omitempty"`    // 用户登录状态
	Roles          []*Role      `json:"roles"`                     // 角色列表
	Projects       []*Project   `json:"projects"`
}

// Password user's password
type Password struct {
	UserID   string `json:"-"`
	Password string `json:"-"`
	ExpireAt int64  `json:"expire_at"`           // 密码过期时间
	CreateAt int64  `json:"create_at"`           // 密码创建时间
	UpdateAt int64  `json:"update_at,omitempty"` // 密码更新时间
}

// LoginStatus 用户登录信息统计, 记录标准: hash({user_id}.{applaction_id}.{grant_type}) 为一条记录
type LoginStatus struct {
	IP            string    `json:"ip,omitempty"`             // 用户登录时的IP地址
	Login         int64     `json:"login,omitempty"`          // 用户最近一次退出系统的时间, 用于评估用户使用系统的时长
	Logout        int64     `json:"logout,omitempty"`         // 用户最近一次登录系统的时间
	GrantType     GrantType `json:"grant_type,omitempty"`     // 用户通过哪种授权方式登录的
	Success       int64     `json:"success,omitempty"`        // 用户登录成功的次数, 及用户访问系统的次数
	Failed        int       `json:"failed,omitempty"`         // 用户连续登录失败的次数, 如果登录成功则清零, 用户实现用户多少
	UserID        string    `json:"user_id,omitempty"`        // 用户ID
	ApplicationID string    `json:"application_id,omitempty"` // 用户应用ID
}

// Invitation code
type Invitation struct {
	Code           string   `json:"code"`                     // 邀请码
	Inviter        string   `json:"inviter"`                  // 邀请人
	Invitee        string   `json:"invitee,omitempty"`        // 被邀人
	InviteeDomain  string   `json:"invitee_domain,omitempty"` // 别邀人域ID
	InvitedTime    int64    `json:"invited_time"`             // 邀请时间
	AcceptTime     int64    `json:"accept_time,omitempty"`    // 被邀人接收邀请的时间
	ExpireTime     int64    `json:"expire_time,omitempty"`    // 邀请码过期时间
	InvitationURI  string   `json:"invitation_uri"`           // 邀请URI
	InviteeRoles   []string `json:"invitee_roles"`            // 赋予被邀人的那些角色
	AccessProjects []string `json:"access_project_ids"`       // 赋予被邀人项目访问范围
}

// Validate 校验创建时的参数
func (u *User) Validate() error {
	if u.Account == "" {
		return exception.NewBadRequest("the user's account required!")
	}

	if len(u.Account) > 128 {
		return exception.NewBadRequest("user's account is too long,  max length is 128")
	}

	if u.Password != nil && len(u.Password.Password) < 6 {
		return exception.NewBadRequest("user password length must not be less than 6")
	}

	if u.Domain == nil || u.Domain.ID == "" {
		return exception.NewBadRequest("the user's domain required!")
	}

	if u.Department == nil || u.Department.ID == "" {
		return exception.NewBadRequest("the user's department required!")
	}

	return nil
}

func (u *User) String() string {
	str, err := json.Marshal(u)
	if err != nil {
		log.Printf("E! marshal user to string error: %s", err)
		return fmt.Sprintf("ID: %s, Name: %s", u.ID, u.Account)
	}

	return string(str)
}

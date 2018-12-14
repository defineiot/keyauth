package user

import (
	"github.com/defineiot/keyauth/dao/domain"
	"github.com/defineiot/keyauth/dao/project"
	"github.com/defineiot/keyauth/internal/exception"
)

// User info
type User struct {
	ID                string    `json:"id"`                  // 用户UUID
	Account           string    `json:"account"`             // 用户账号名称
	Mobile            string    `json:"mobile"`              // 手机号码, 用户可以通过手机进行注册和密码找回, 还可以通过手机号进行登录
	Email             string    `json:"email"`               // 邮箱, 用户可以通过邮箱进行注册和照明密码
	Phone             string    `json:"phone"`               // 用户的座机号码
	Address           string    `json:"address"`             // 用户住址
	RealName          string    `json:"real_name"`           // 用户真实姓名
	NickName          string    `json:"nick_name"`           // 用户昵称, 用于在界面进行展示
	Gender            string    `json:"gender"`              // 性别
	Avatar            string    `json:"avatar"`              // 头像
	Locked            string    `json:"locked"`              // 是否冻结次用户
	LastLogoutTime    string    `json:"last_logout_time"`    // 用户最近一次退出系统的时间, 用于评估用户使用系统的时长
	LastLoginTime     string    `json:"last_login_time"`     // 用户最近一次登录系统的时间
	LastLoginIP       string    `json:"last_login_ip"`       // 用户最近一次登录系统的客户端IP地址
	LoginFailedTimes  string    `json:"login_failed_times"`  // 用户连续登录失败的次数, 如果登录成功则清零
	LoginSuccessTimes string    `json:"login_success_times"` // 用户登录成功的次数, 及用户访问系统的次数
	CreateAt          string    `json:"create_at"`           // 用户创建的时间
	ExpiresActiveDays string    `json:"expires_active_days"` // 用户多久未登录时(天), 冻结改用户, 防止僵尸用户的账号被利用
	Password          *Password `json:"password"`            // 密码相关信息

	Domain         *domain.Domain   `json:"domain"`     // 如果需要对象由上层进行查找
	DefaultProject *project.Project `json:"project"`    //  如果需要对象由上层进行查找
	Department     *Department      `json:"department"` // 所属部门信息
	RoleNames      []string         `json:"roles"`      // 角色列表
}

// Password user's password
type Password struct {
	ID       int64  `json:"-"`
	UserID   string `json:"-"`
	Password string `json:"-"`
	ExpireAt int64  `json:"expire_at"`           // 密码过期时间
	CreateAt int64  `json:"create_at"`           // 密码创建时间
	UpdateAt int64  `json:"update_at,omitempty"` // 密码更新时间
}

// Validate 校验创建时的参数
func (u *User) Validate() error {
	if u.Account == "" {
		return exception.NewBadRequest("the user's account required!")
	}

	if len(u.Account) > 128 {
		return exception.NewBadRequest("user's account is too long,  max length is 128")
	}

	if u.Password == nil {
		return exception.NewBadRequest("the user's password required!")
	}

	if len(u.Password.Password) < 6 {
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

// Store is user service
type Store interface {
	Reader
	Writer
	Close() error
}

// Reader use to read user information form store
type Reader interface {
	ListUser(domainID string) ([]*User, error)
	ListUserRoles(domainID, userID string) ([]string, error)
	GetUserByID(userID string) (*User, error)
	GetUserByName(domainID, userName string) (*User, error)
	ValidateUser(domainID, userName, password string) (string, error)
	ValidateGlobalUser(userName, password string) (string, error)
	CheckUserNameIsExist(domainID, userName string) (bool, error)
	CheckUserNameIsGlobalExist(userName string) (bool, error)
	CheckUserIsExistByID(userID string) (bool, error)

	ListUserProjects(domainID, userID string) ([]string, error)
	ListUserOtherDomains(userID string) ([]string, error)
}

// Writer use to write user information to store
type Writer interface {
	CreateUser(u *User) (*User, error)

	SaveVerifyCode(toEmail, phoneNumber string, code int) (*VerifyCode, error)
	GetVerifyCodeByMail(toEmail string, code int) (*VerifyCode, error)
	GetVerifyCodeByPhone(phoneNumber string, code int) (*VerifyCode, error)
	RevolkVerifyCode(id int64) error

	SaveUserOtherDomain(userID, otherDomainID string) error
	DeleteUserOtherDomain(userID, otherDomainID string) error

	DeleteUser(domainID, userID string) error

	SetUserPassword(userID, oldPass, newPass string) error
	SetDefaultProject(domainID, userID, projectID string) error
	AddProjectsToUser(domainID, userID string, projectIDs ...string) error
	RemoveProjectsFromUser(domainID, userID string, projectIDs ...string) error
	BindRole(domainID, userID, roleName string) error
	UnBindRole(domainID, userID, roleName string) error

	SaveInvitationsRecord(inviterID string, invitedRoles, accessProjects []string) (*InvitationCode, error)
	ListInvitationRecord(inviterID string) ([]*InvitationCode, error)
	GetInvitationRecord(inviterID, code string) (*InvitationCode, error)
	DeleteInvitationRecord(id int64) error
	UpdateInvitationsRecord(ir *InvitationCode) error
}

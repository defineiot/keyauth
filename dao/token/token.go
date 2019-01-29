package token

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/defineiot/keyauth/dao/role"

	"github.com/defineiot/keyauth/dao/project"

	"github.com/defineiot/keyauth/internal/exception"
)

// GrantType is the type for OAuth2 param `grant_type`
type GrantType string

// Type token type
type Type string

const (
	// oauth2 Authorization Grant: https://tools.ietf.org/html/rfc6749#section-1.3

	AUTHCODE GrantType = "authorization_code" // AUTHCODE oauth2 Authorization Code Grant
	IMPLICIT GrantType = "implicit"           // IMPLICIT oauth2 Implicit Grant
	PASSWORD GrantType = "password"           // PASSWORD oauth2 Resource Owner Password Credentials Grant
	CLIENT   GrantType = "client_credentials" // CLIENT oauth2 Client Credentials Grant
	REFRESH  GrantType = "refresh_token"      // REFRESH oauth2 Refreshing an Access Token
	UPSCOPE  GrantType = "upgrade_scope"      // UPSCOPE is an custom grant for use unscope token acquire scope token
	WeChat   GrantType = "wechat"             // WeChat is an custom grant for use unscope token acquire scope token

	// oauth2 Token Type: https://tools.ietf.org/html/rfc6749#section-7.1

	Bearer Type = "bearer" // detail: https://tools.ietf.org/html/rfc6750
	MAC    Type = "mac"    // detail: https://tools.ietf.org/html/rfc6749#ref-OAuth-HTTP-MAC
	JWT    Type = "jwt"    // detail: https://tools.ietf.org/html/rfc7519
)

// Code is oauth2 auth code https://tools.ietf.org/html/rfc6749#section-4.1.2
type Code struct {
	Code  string
	State string
}

// Token is user's access resource token
type Token struct {
	AccessToken    string    `json:"access_token"`              // 服务访问令牌
	RefreshToken   string    `json:"refresh_token,omitempty"`   // 用于刷新访问令牌的凭证, 刷新过后, 原先令牌将会被删除
	TokenType      Type      `json:"token_type,omitempty"`      // 令牌的类型 类型包含: bearer/jwt  (默认为bearer)
	GrantType      GrantType `json:"grant_type,omitempty"`      // 授权的类型
	UserID         string    `json:"user_id,omitempty"`         // 用户ID
	CurrentProject string    `json:"current_project,omitempty"` // 当前所在项目
	DomainID       string    `json:"domain_id,omitempty"`       // 用户所在的域的ID, 用户可以切换域(如果用户加入了多个域)
	ServiceID      string    `json:"service_id,omitempty"`      // 服务ID, 如果凭证是颁发给内部服务使用时, 服务删除时,颁发给它的令牌需要删除, 服务禁用时, 令牌验证不通过
	ApplicationID  string    `json:"application_id,omitempty"`  // 用户应用ID, 如果凭证是颁发给应用的, 应用在删除时需要删除所有的令牌, 应用禁用时, 该应用令牌验证会不通过
	Name           string    `json:"name,omitempty"`            // 独立颁发给SDK使用时, 命名token
	Description    string    `json:"description,omitempty"`     // 独立颁发给SDK使用时, 令牌的描述信息, 方便定位与取消
	Scope          string    `json:"scope,omitempty"`           // 令牌的作用范围: detail https://tools.ietf.org/html/rfc6749#section-3.3
	CreatedAt      int64     `json:"create_at,omitempty"`       // 凭证创建时间
	ExpiresIn      int64     `json:"expires_in,omitempty"`      // 凭证过期的时间

	IsSystemAdmin     bool               `json:"is_system_admin,omitempty"`    // 是否是系统管理员
	IsDomainAdmin     bool               `json:"is_domain_admin,omitempty"`    // 是否是域管理员
	Roles             []*role.Role       `json:"roles,omitempty"`              // 该凭证的权限列表
	AvaliableProjects []*project.Project `json:"available_projects,omitempty"` // 该用户可以访问的项目列表
}

// Store is auth service
type Store interface {
	StoreReader
	StoreWriter
	Close() error
}

// StoreReader read information from store
type StoreReader interface {
	GetToken(accessToken string) (*Token, error)
	GetTokenByRefresh(refreshToken string) (*Token, error)
}

// StoreWriter write information to store
type StoreWriter interface {
	SaveToken(t *Token) error
	DeleteTokenByRefresh(refreshToken string) error
	UpdateTokenScope(accessToken, scope string) error
	DeleteToken(accessToken string) error
}

func (t *Token) String() string {
	str, err := json.Marshal(t)
	if err != nil {
		log.Printf("E! marshal role to string error: %s", err)
		return fmt.Sprintf("access_token: %s, refresh_token: %s", t.AccessToken, t.RefreshToken)
	}

	return string(str)
}

// ValidateSave 校验token创建
func (t *Token) ValidateSave() error {
	if t.UserID == "" && t.ServiceID == "" {
		return exception.NewBadRequest("token's user_id or service_id is missed")
	}
	if t.ServiceID == "" && t.ApplicationID == "" {
		return exception.NewBadRequest("token's service_id or application_id required!")
	}
	if t.AccessToken == "" {
		return exception.NewInternalServerError("token's access token must'nt be \"\"")
	}
	if t.TokenType == "" {
		return exception.NewInternalServerError("token's type must one of bearer or jwt")
	}
	if t.GrantType != AUTHCODE && t.GrantType != IMPLICIT && t.GrantType != PASSWORD && t.GrantType != CLIENT && t.GrantType != REFRESH {
		return exception.NewBadRequest("grant_type must one of authorization_code,implicit,password,client_credentials,refresh_token")
	}

	return nil
}

// IsExpired use to validate the token is expired
func (t *Token) IsExpired() (bool, int64) {
	now := time.Now().Unix()
	allow := t.CreatedAt + t.ExpiresIn

	if now < allow {
		return true, allow - now
	}

	return false, 0
}

// MakeBearerToken https://tools.ietf.org/html/rfc6750#section-2.1
// b64token    = 1*( ALPHA / DIGIT /"-" / "." / "_" / "~" / "+" / "/" ) *"="
func MakeBearerToken(lenth int) string {
	charlist := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-._~+/"
	t := make([]string, lenth)
	rand.Seed(time.Now().UnixNano() + int64(lenth) + rand.Int63n(10000))
	for i := 0; i < lenth; i++ {
		rn := rand.Intn(len(charlist))
		w := charlist[rn : rn+1]
		t = append(t, w)
	}

	token := strings.Join(t, "")
	return base64.RawURLEncoding.EncodeToString([]byte(token))
}

package store

import (
	"encoding/base64"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/defineiot/keyauth/dao/token"
	"github.com/defineiot/keyauth/dao/user"
	"github.com/defineiot/keyauth/internal/exception"
)

const (
	PKCE_PLAIN = "plain" // PKCE_PLAIN is oauth pkce extension
	PKCE_S256  = "S256"  // PKCE_S256 is oauth pkce extension
)

var (
	pkceMatcher = regexp.MustCompile("^[a-zA-Z0-9~._-]{43,128}$")
)

// TokenRequest use to request access token
type TokenRequest struct {
	GrantType           token.GrantType `json:"grant_type,omitempty"`
	ClientID            string          `json:"client_id,omitempty"`
	ClientSecret        string          `json:"client_secret,omitempty"`
	AuthorizationHeader string          `json:"authorization_header,omitempty"`
	DomainID            string          `json:"domain_id,omitempty"`
	Username            string          `json:"username,omitempty"`
	Password            string          `json:"password,omitempty"`
	Code                string          `json:"code,omitempty"`
	RedirectURI         string          `json:"redirect_uri,omitempty"`
	RefreshToken        string          `json:"refresh_token,omitempty"`
	AccessToken         string          `json:"access_token,omitempty"`
	Scope               string          `json:"scope,omitempty"`

	isCheckClient bool
}

//  validate the request
func (t *TokenRequest) validate() error {
	switch t.GrantType {
	case token.AUTHCODE:
		if t.Code == "" {
			return exception.NewBadRequest("if use %s grant type, code is needed", t.GrantType)
		}
		if t.RedirectURI == "" {
			return exception.NewBadRequest("if use %s grant type, redirect uri is need", t.GrantType)
		}
		t.isCheckClient = true

	case token.IMPLICIT:
		if t.ClientID == "" {
			return exception.NewBadRequest("if use %s grant type, client id is needed", t.GrantType)
		}
		if t.RedirectURI == "" {
			return exception.NewBadRequest("if use %s grant type, redirect uri is need", t.GrantType)
		}

	case token.PASSWORD:
		if t.Username == "" || t.Password == "" {
			return exception.NewBadRequest("if use %s grant type, username, password is needed", t.GrantType)
		}
		t.isCheckClient = true

	case token.CLIENT:
		t.isCheckClient = true

	case token.REFRESH:
		if t.RefreshToken == "" {
			return exception.NewBadRequest("if use %s grant type, refresh token is needed", t.GrantType)
		}
		t.isCheckClient = true

	case token.UPSCOPE:
		if t.AccessToken == "" {
			return exception.NewBadRequest("if use %s grant type, access_token is needed", t.GrantType)
		}
		if t.Scope == "" {
			return exception.NewBadRequest("if use %s grant type, scope project_id is needed", t.GrantType)
		}
		t.isCheckClient = true

	default:
		return exception.NewBadRequest(`invalid_grant, supported grant type: [authorization_code, 
		implicit, password, client_credentials, refresh_token, upgrade_scope]`)
	}

	if t.isCheckClient {
		if t.ClientID == "" || t.ClientSecret == "" {
			return exception.NewBadRequest("if use %s grant type, client id and client secret is needed",
				t.GrantType)
		}
	}

	return nil
}

// IssueToken use to issue access token
func (s *Store) IssueToken(req *TokenRequest) (*token.Token, error) {
	if err := req.validate(); err != nil {
		return nil, err
	}

	var t *token.Token
	switch req.GrantType {
	case token.AUTHCODE:
		app, err := s.dao.Application.GetApplicationByClientID(req.ClientID)
		if err != nil {
			return nil, err
		}
		if req.ClientSecret != app.ClientSecret {
			return nil, exception.NewUnauthorized("unauthorized_client")
		}

		t, err = s.issueTokenByAuthCode(app.ID, req.Code, req.RedirectURI)
		if err != nil {
			return nil, err
		}

	case token.IMPLICIT:
		app, err := s.dao.Application.GetApplicationByClientID(req.ClientID)
		if err != nil {
			return nil, err
		}

		t, err = s.issueTokenByImplicit(app.ID, req.RedirectURI)
		if err != nil {
			return nil, err
		}

	case token.PASSWORD:
		app, err := s.dao.Application.GetApplicationByClientID(req.ClientID)
		if err != nil {
			return nil, err
		}
		if req.ClientSecret != app.ClientSecret {
			return nil, exception.NewUnauthorized("unauthorized_client")
		}

		t, err = s.issueTokenByPassword(req.Scope, app.ID, req.Username, req.Password)
		if err != nil {
			return nil, err
		}

	case token.CLIENT:
		svr, err := s.dao.Service.GetServiceByClientID(req.ClientID)
		if err != nil {
			return nil, err
		}
		if req.ClientSecret != svr.ClientSecret {
			return nil, exception.NewUnauthorized("unauthorized_client")
		}

		t, err = s.issueTokenByClient(svr.ID, req.Scope)
		if err != nil {
			return nil, err
		}

	case token.REFRESH:
		app, err := s.dao.Application.GetApplicationByClientID(req.ClientID)
		if err != nil {
			return nil, err
		}
		if req.ClientSecret != app.ClientSecret {
			return nil, exception.NewUnauthorized("unauthorized_client")
		}

		t, err = s.issueTokenByRefresh(req.RefreshToken)
		if err != nil {
			return nil, err
		}

	case token.UPSCOPE:
		app, err := s.dao.Application.GetApplicationByClientID(req.ClientID)
		if err != nil {
			return nil, err
		}
		if req.ClientSecret != app.ClientSecret {
			return nil, exception.NewUnauthorized("unauthorized_client")
		}

		t, err = s.issueTokenByUpScope(req.AccessToken, req.Scope)
		if err != nil {
			return nil, err
		}

	default:
		return nil, exception.NewBadRequest(`invalid_grant only support 
		[authorization_code, implicit, password, client_credentials, refresh_token, upgrade_scope]`)
	}

	return t, nil
}

// ValidateToken use to validate token
func (s *Store) ValidateToken(accessToken string) (*token.Token, error) {
	var (
		err    error
		cached bool
	)

	tk := new(token.Token)
	// 尝试从缓存中获取Token
	cacheKey := s.cachePrefix.token + accessToken
	if s.isCache {
		if s.cache.Get(cacheKey, tk) {
			s.log.Debug("get token from cache key: %s", cacheKey)
			cached = true
		}
		s.log.Debug("get token from cache failed, key: %s", cacheKey)
	}

	if !cached {
		// 如果没有从缓存中获取到, 则从到层获取
		tk, err = s.dao.Token.GetToken(accessToken)
		if err != nil {
			return nil, err
		}
	}

	// 计算token是否过期
	ok, delta := tk.IsExpired()
	if !ok {
		return nil, exception.NewExpired("token has expired, access_token: %s", tk.AccessToken)
	}
	tk.ExpiresIn = delta

	// 获取用户可以访问的项目列表
	if tk.UserID != "" {
		projects, err := s.dao.Project.ListUserProjects(tk.DomainID, tk.UserID)
		if err != nil {
			return nil, exception.NewInternalServerError(err.Error())
		}

		tk.AvaliableProjects = projects
	}

	// 获取用户的角色列表
	roles, err := s.dao.Role.ListUserRole(tk.DomainID, tk.UserID)
	if err != nil {
		return nil, err
	}
	tk.Roles = roles

	// 更新缓存
	if s.isCache {
		if !s.cache.Set(cacheKey, tk, s.ttl) {
			s.log.Debug("set token cache failed, key: %s", cacheKey)
		}
		s.log.Debug("set token cache ok, key: %s", cacheKey)
	}

	return tk, nil
}

// RevokeToken refresh token
func (s *Store) RevokeToken(accessToken string) error {
	var err error

	cacheKey := s.cachePrefix.token + accessToken

	err = s.dao.Token.DeleteToken(accessToken)
	if err != nil {
		return err
	}

	if s.isCache {
		if !s.cache.Delete(cacheKey) {
			s.log.Debug("delete token from cache failed, key: %s", cacheKey)
		}
		s.log.Debug("delete token from cache success, key: %s", cacheKey)
	}

	return nil
}

// issueTokenByAuthCode implement Authorization Code Grant
// https://tools.ietf.org/html/rfc6749#section-4.1.3
func (s *Store) issueTokenByAuthCode(clientID, code, redirectURI string) (*token.Token, error) {
	return nil, nil
}

// issueTokenByImplicit implement Implicit Grant
// https://tools.ietf.org/html/rfc6749#section-4.2
func (s *Store) issueTokenByImplicit(clientID, redirectURI string) (*token.Token, error) {
	return nil, nil
}

// issueTokenByPassword implement Resource Owner Password Credentials Grant
// https://tools.ietf.org/html/rfc6749#section-4.3
func (s *Store) issueTokenByPassword(scope, appID, account, password string) (*token.Token, error) {
	// 查询用户是否存在
	user, err := s.dao.User.GetUser(user.Account, account)
	if err != nil {
		return nil, err
	}

	// 检查用户的秘密是否正确
	if s.hmacHash(password) != user.Password.Password {
		return nil, exception.NewForbidden("username or password invalidate")
	}

	// 生成Token
	tk, err := s.generateToken(scope, user.Domain.ID, user.ID, appID, token.Bearer, token.PASSWORD)
	if err != nil {
		return nil, err
	}

	// 获取token的项目列表
	projects, err := s.dao.Project.ListUserProjects(user.Domain.ID, user.ID)
	if err != nil {
		return nil, exception.NewInternalServerError(err.Error())
	}
	tk.AvaliableProjects = projects

	// 获取用户的角色
	roles, err := s.dao.Role.ListUserRole(user.Domain.ID, user.ID)
	if err != nil {
		return nil, err
	}
	tk.Roles = roles

	for i := range roles {
		switch roles[i].Name {
		case "system_admin":
			tk.IsSystemAdmin = true
		case "domain_admin":
			tk.IsDomainAdmin = true
		}
	}

	return tk, nil
}

// issueTokenByClient implement Client Credentials Grant
// https://tools.ietf.org/html/rfc6749#section-4.4.2
func (s *Store) issueTokenByClient(serviceID string, scope string) (*token.Token, error) {
	t := new(token.Token)
	t.Scope = scope
	t.CreatedAt = time.Now().Unix()
	t.ExpiresIn = s.conf.Token.ExpiresIn
	t.GrantType = token.CLIENT
	t.ServiceID = serviceID

	switch t.TokenType {
	case "bearer", "":
		t.TokenType = token.Bearer
		t.AccessToken = makeBearerToken(24)
		t.RefreshToken = makeBearerToken(32)
	case "jwt":
		t.TokenType = token.JWT
	default:
		return nil, exception.NewInternalServerError("unknown token type, %s", t.TokenType)
	}

	if err := s.dao.Token.SaveToken(t); err != nil {
		return nil, err
	}

	return t, nil
}

// issueTokenByRefresh implement Refreshing an Access Token
// https://tools.ietf.org/html/rfc6749#section-6
func (s *Store) issueTokenByRefresh(refreshToken string) (*token.Token, error) {
	if refreshToken == "" {
		return nil, exception.NewBadRequest("resfresh_token missed")
	}

	// 通过refresh_token找到用户的token
	old, err := s.dao.Token.GetTokenByRefresh(refreshToken)
	if err != nil {
		return nil, err
	}

	// 删除旧的token
	if err := s.dao.Token.DeleteToken(old.AccessToken); err != nil {
		return nil, err
	}

	// 清除缓存的token
	if s.isCache {
		cacheKey := s.cachePrefix.token + old.AccessToken
		if !s.cache.Delete(cacheKey) {
			s.log.Debug("delete token from cache failed, key: %s", cacheKey)
		}
		s.log.Debug("delete token from cache success, key: %s", cacheKey)
	}

	// 生成新token
	tk, err := s.generateToken(old.Scope, old.DomainID, old.UserID, old.ApplicationID, token.Bearer, token.REFRESH)
	if err != nil {
		return nil, err
	}

	// 新token项目不变
	tk.AvaliableProjects = old.AvaliableProjects
	tk.Roles = old.Roles
	tk.IsSystemAdmin = old.IsSystemAdmin
	tk.IsDomainAdmin = old.IsDomainAdmin

	return tk, nil
}

func (s *Store) issueTokenByUpScope(accessToken, scope string) (*token.Token, error) {
	if accessToken == "" {
		return nil, exception.NewBadRequest("access_token missed")
	}

	scopeSlice := strings.Split(scope, ",")
	if len(scopeSlice) != 2 {
		return nil, exception.NewBadRequest("scope format invalidate, format: <domain_id>,<project_id>")
	}

	domainID, projectID := scopeSlice[0], scopeSlice[1]
	if projectID == "" && domainID == "" {
		return nil, exception.NewBadRequest("scope project_id or domain_id missed")
	}

	// 校验当前Token是否合法
	t, err := s.ValidateToken(accessToken)
	if err != nil {
		return nil, err
	}

	// 切换用户的域空间, 判断需要切换的域是否属于该用户
	if domainID != "" && domainID != t.DomainID {
		var validateD bool
		otherDs, err := s.dao.Domain.ListUserThirdDomains(t.UserID)
		if err != nil {
			return nil, err
		}
		for _, ad := range otherDs {
			if domainID == ad.ID {
				validateD = true
				break
			}
		}
		if !validateD {
			return nil, exception.NewForbidden("the domain: %s not belong user: %s", domainID, t.UserID)
		}
	}

	// 切换用户的项目空间, 判断需要切换的项目是否属于该用户
	var projectOK bool
	for i := range t.AvaliableProjects {
		if t.AvaliableProjects[i].ID == projectID {
			projectOK = true
			break
		}
	}
	if !projectOK {
		return nil, exception.NewBadRequest("the project: %s not belong user: %s", projectID, t.UserID)
	}

	// 更新Token的Scope
	if err := s.dao.Token.UpdateTokenScope(t.AccessToken, scope); err != nil {
		return nil, err
	}
	t.Scope = scope

	return t, nil
}

// https://tools.ietf.org/html/rfc6750#section-2.1
// b64token    = 1*( ALPHA / DIGIT /"-" / "." / "_" / "~" / "+" / "/" ) *"="
func makeBearerToken(lenth int) string {
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

func (s *Store) generateToken(scope, domainID, userID, appID string, tp token.Type, gt token.GrantType) (*token.Token, error) {
	t := &token.Token{
		Scope:         scope,
		DomainID:      domainID,
		CreatedAt:     time.Now().Unix(),
		ExpiresIn:     s.conf.Token.ExpiresIn,
		TokenType:     tp,
		UserID:        userID,
		ApplicationID: appID,
		GrantType:     gt,
	}

	switch tp {
	case token.Bearer:
		t.AccessToken = makeBearerToken(24)
		t.RefreshToken = makeBearerToken(32)
	case token.JWT:
	default:
		return nil, exception.NewInternalServerError("unknown token type, %s, only support bearer", tp)
	}

	if err := s.dao.Token.SaveToken(t); err != nil {
		return nil, err
	}

	return t, nil
}

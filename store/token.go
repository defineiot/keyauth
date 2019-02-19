package store

import (
	"encoding/base64"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/defineiot/keyauth/dao/models"
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
	GrantType           models.GrantType `json:"grant_type,omitempty"`
	ClientID            string           `json:"client_id,omitempty"`
	ClientSecret        string           `json:"client_secret,omitempty"`
	AuthorizationHeader string           `json:"authorization_header,omitempty"`
	DomainID            string           `json:"domain_id,omitempty"`
	Username            string           `json:"username,omitempty"`
	Password            string           `json:"password,omitempty"`
	Code                string           `json:"code,omitempty"`
	RedirectURI         string           `json:"redirect_uri,omitempty"`
	RefreshToken        string           `json:"refresh_token,omitempty"`
	AccessToken         string           `json:"access_token,omitempty"`
	Scope               string           `json:"scope,omitempty"`

	isCheckClient bool
}

//  validate the request
func (t *TokenRequest) validate() error {
	switch t.GrantType {
	case models.AUTHCODE:
		if t.Code == "" {
			return exception.NewBadRequest("if use %s grant type, code is needed", t.GrantType)
		}
		if t.RedirectURI == "" {
			return exception.NewBadRequest("if use %s grant type, redirect uri is need", t.GrantType)
		}
		t.isCheckClient = true

	case models.IMPLICIT:
		if t.ClientID == "" {
			return exception.NewBadRequest("if use %s grant type, client id is needed", t.GrantType)
		}
		if t.RedirectURI == "" {
			return exception.NewBadRequest("if use %s grant type, redirect uri is need", t.GrantType)
		}

	case models.PASSWORD:
		if t.Username == "" || t.Password == "" {
			return exception.NewBadRequest("if use %s grant type, username, password is needed", t.GrantType)
		}
		t.isCheckClient = true

	case models.CLIENT:
		t.isCheckClient = true

	case models.REFRESH:
		if t.RefreshToken == "" {
			return exception.NewBadRequest("if use %s grant type, refresh token is needed", t.GrantType)
		}
		t.isCheckClient = true

	case models.UPSCOPE:
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
func (s *Store) IssueToken(req *TokenRequest) (t *models.Token, err error) {
	if err := req.validate(); err != nil {
		return nil, err
	}

	switch req.GrantType {
	case models.AUTHCODE:
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

	case models.IMPLICIT:
		app, err := s.dao.Application.GetApplicationByClientID(req.ClientID)
		if err != nil {
			return nil, err
		}

		t, err = s.issueTokenByImplicit(app.ID, req.RedirectURI)
		if err != nil {
			return nil, err
		}

	case models.PASSWORD:
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

	case models.CLIENT:
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

	case models.REFRESH:
		t, err = s.issueTokenByRefresh(req.ClientID, req.ClientSecret, req.RefreshToken)
		if err != nil {
			return nil, err
		}

	case models.UPSCOPE:
		t, err = s.issueTokenByUpScope(req.ClientID, req.ClientSecret, req.AccessToken, req.Scope)
		if err != nil {
			return nil, err
		}

	default:
		return nil, exception.NewBadRequest(`invalid_grant only support 
		[authorization_code, implicit, password, client_credentials, refresh_token, upgrade_scope]`)
	}

	return t, nil
}

// ValidateTokenReq token校验
type ValidateTokenReq struct {
	ClientID     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	FeatureName  string `json:"feature_name,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

func (v *ValidateTokenReq) validate() error {
	if v.ClientID == "" || v.ClientSecret == "" {
		return exception.NewForbidden("client credentials missed")
	}

	if v.AccessToken == "" {
		return exception.NewBadRequest("token missed")
	}

	return nil
}

// ValidateTokenWithClient use to validate token
func (s *Store) ValidateTokenWithClient(v *ValidateTokenReq) (*models.Token, error) {
	var (
		err    error
		cached bool
	)

	// 校验请求的合法性
	if err := v.validate(); err != nil {
		return nil, err
	}

	// 校验后端服务调用的合法性
	svr, err := s.dao.Service.GetServiceByClientID(v.ClientID)
	if err != nil {
		s.log.Debug("find service by client error, %s", err)
	}
	if svr != nil {
		if v.ClientSecret != svr.ClientSecret {
			return nil, exception.NewForbidden("unauthorized_client")
		}
	}

	// 校验前段应用调用的合法性
	app, err := s.dao.Application.GetApplicationByClientID(v.ClientID)
	if err != nil {
		s.log.Debug("find application by client error, %s", err)
	}
	if app != nil {
		if v.ClientSecret != app.ClientSecret {
			return nil, exception.NewForbidden("unauthorized_client")
		}
	}

	tk := new(models.Token)
	// 尝试从缓存中获取Token
	cacheKey := s.cachePrefix.token + v.AccessToken
	if s.isCache {
		if s.cache.Get(cacheKey, tk) {
			s.log.Debug("get token from cache key: %s", cacheKey)
			cached = true
		} else {
			s.log.Debug("get token from cache failed, key: %s", cacheKey)
		}
	}

	if !cached {
		// 如果没有从缓存中获取到, 则从DAO层获取
		tk, err = s.dao.Token.GetToken(v.AccessToken)
		if err != nil {
			return nil, err
		}

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
	}

	// 计算token是否过期
	ok, delta := tk.IsExpired()
	if !ok {
		return nil, exception.NewExpired("token has expired, access_token: %s", tk.AccessToken)
	}
	tk.ExpiresAt = delta

	// 应用调用者只能校验自己的token
	if app != nil {
		if tk.ApplicationID != app.ID {
			return nil, exception.NewForbidden("this is not your token")
		}
	}

	// 服务调用者可以校验其他人的token

	// 校验用户是否有权限访问指定的功能
	if v.FeatureName != "" || svr != nil {
		switch tk.GrantType {
		case models.CLIENT:
			if v.FeatureName != "RegistryServiceFeatures" {
				return nil, exception.NewForbidden("client_credentials only can acess RegistryServiceFeatures")
			}
		case models.PASSWORD, models.UPSCOPE, models.REFRESH:
			if v.FeatureName == "RegistryServiceFeatures" {
				return nil, exception.NewForbidden("client_credentials only can acess RegistryServiceFeatures")
			}
		default:
			return nil, exception.NewBadRequest("other grant type not support")
		}

		var hasPerm bool
		for i := range tk.Roles {
			rn := tk.Roles[i].Name
			if rn == systemAdminRoleName {
				hasPerm = true
				tk.IsSystemAdmin = true
				continue
			}

			if rn == domainAdminRoleName {
				hasPerm = true
				tk.IsDomainAdmin = true
				break
			}

			role, err := s.GetRole(rn)
			if err != nil {
				return nil, err
			}

			for _, f := range role.Features {
				if f.Name == v.FeatureName {
					hasPerm = true
				}
			}
		}

		if !hasPerm {
			return nil, exception.NewForbidden("user: %s has no permisson for access feature: %s", tk.UserID, v.FeatureName)
		}
	}

	// 更新缓存
	if s.isCache {
		if !s.cache.Set(cacheKey, tk, s.ttl) {
			s.log.Debug("set token cache failed, key: %s", cacheKey)
		}
		s.log.Debug("set token cache ok, key: %s", cacheKey)
	}

	tk.RefreshToken = ""
	return tk, nil
}

// ValidateToken use to validate token
func (s *Store) ValidateToken(accessToken, featureName string) (*models.Token, error) {
	var (
		err    error
		cached bool
	)

	tk := new(models.Token)
	// 尝试从缓存中获取Token
	cacheKey := s.cachePrefix.token + accessToken
	if s.isCache {
		if s.cache.Get(cacheKey, tk) {
			s.log.Debug("get token from cache key: %s", cacheKey)
			cached = true
		} else {
			s.log.Debug("get token from cache failed, key: %s", cacheKey)
		}
	}

	if !cached {
		// 如果没有从缓存中获取到, 则从DAO层获取
		tk, err = s.dao.Token.GetToken(accessToken)
		if err != nil {
			return nil, err
		}

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
	}

	// 计算token是否过期
	ok, delta := tk.IsExpired()
	if !ok {
		return nil, exception.NewExpired("token has expired, access_token: %s", tk.AccessToken)
	}
	tk.ExpiresAt = delta

	// 校验用户是否有权限访问指定的功能
	if featureName != "" {
		var hasPerm bool

		switch tk.GrantType {
		case models.CLIENT:
			if featureName != "RegistryServiceFeatures" {
				return nil, exception.NewForbidden("client_credentials only can acess RegistryServiceFeatures")
			}
			hasPerm = true
		case models.PASSWORD, models.UPSCOPE, models.REFRESH:
			if featureName == "RegistryServiceFeatures" {
				return nil, exception.NewForbidden("client_credentials only can acess RegistryServiceFeatures")
			}
		default:
			return nil, exception.NewBadRequest("other grant type not support")
		}

		for i := range tk.Roles {
			rn := tk.Roles[i].Name
			if rn == systemAdminRoleName {
				hasPerm = true
				tk.IsSystemAdmin = true
				continue
			}

			if rn == domainAdminRoleName {
				hasPerm = true
				tk.IsDomainAdmin = true
				break
			}

			role, err := s.GetRole(rn)
			if err != nil {
				return nil, err
			}

			for _, f := range role.Features {
				if f.Name == featureName {
					hasPerm = true
				}
			}
		}

		if !hasPerm {
			return nil, exception.NewForbidden("user: %s has no permisson for access feature: %s", tk.UserID, featureName)
		}
	}

	// 更新缓存
	if s.isCache {
		if !s.cache.Set(cacheKey, tk, s.ttl) {
			s.log.Debug("set token cache failed, key: %s", cacheKey)
		}
		s.log.Debug("set token cache ok, key: %s", cacheKey)
	}

	return tk, nil
}

// RevokeReq 撤销Token
type RevokeReq struct {
	AccessToken  string `json:"access_token,omitempty"`
	ClientID     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
}

func (r *RevokeReq) validate() error {
	if r.AccessToken == "" {
		return exception.NewBadRequest("access_token missed")
	}

	if r.ClientID == "" || r.ClientSecret == "" {
		return exception.NewForbidden("client credentials missed")
	}

	return nil
}

// RevokeToken refresh token
func (s *Store) RevokeToken(req *RevokeReq) error {
	var err error

	// 校验请求的合法性
	if err := req.validate(); err != nil {
		return err
	}

	// 获取要被撤销的token对象
	tk, err := s.dao.Token.GetToken(req.AccessToken)
	if err != nil {
		return err
	}

	// 校验后端服务调用的合法性(服务调用者可以校验其他人的token)
	svr, err := s.dao.Service.GetServiceByClientID(req.ClientID)
	if err != nil {
		s.log.Debug("find service by client error, %s", err)
	}
	if svr != nil {
		if req.ClientSecret != svr.ClientSecret {
			return exception.NewForbidden("unauthorized_client")
		}

		if tk.ServiceID != svr.ID {
			return exception.NewForbidden("this is not your token")
		}
	}

	// 校验前段应用调用的合法性(应用调用者 只能校验自己的token)
	app, err := s.dao.Application.GetApplicationByClientID(req.ClientID)
	if err != nil {
		s.log.Debug("find application by client error, %s", err)
	}
	if app != nil {
		if req.ClientSecret != app.ClientSecret {
			return exception.NewForbidden("unauthorized_client")
		}

		if tk.ApplicationID != app.ID {
			return exception.NewForbidden("this is not your token")
		}
	}

	cacheKey := s.cachePrefix.token + req.AccessToken

	err = s.dao.Token.DeleteToken(req.AccessToken)
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
func (s *Store) issueTokenByAuthCode(clientID, code, redirectURI string) (*models.Token, error) {
	return nil, nil
}

// issueTokenByImplicit implement Implicit Grant
// https://tools.ietf.org/html/rfc6749#section-4.2
func (s *Store) issueTokenByImplicit(clientID, redirectURI string) (*models.Token, error) {
	return nil, nil
}

// issueTokenByPassword implement Resource Owner Password Credentials Grant
// https://tools.ietf.org/html/rfc6749#section-4.3
func (s *Store) issueTokenByPassword(scope, appID, account, password string) (*models.Token, error) {
	var tk *models.Token

	// 查询用户是否存在
	user, err := s.dao.User.GetUser(models.AccountIndex, account)
	if err != nil {
		return nil, err
	}

	// 检查用户的密码是否正确
	if s.hmacHash(password) != user.Password.Password {
		return nil, exception.NewForbidden("username or password invalidate")
	}

	// 查看最新的, 还有至少一半时间可用的token, 如果有就使用老的token
	ctk, err := s.dao.Token.GetUserCurrentToken(user.ID, appID, models.PASSWORD)
	if err != nil {
		if _, ok := err.(*exception.NotFound); !ok {
			return nil, err
		}
	}

	if ctk != nil {
		if ok, delta := ctk.IsExpired(); ok && delta > ctk.ExpiresIn/2 {
			tk = ctk
		} else {
			// 如果token所剩时间不足一半, 则清除该token
			if err := s.dao.Token.DeleteToken(ctk.AccessToken); err != nil {
				s.log.Warn("clean expired token error, %s", err)
			}
		}
	}

	// 生成新Token
	if tk == nil {
		tk, err = s.generateToken(scope, user.Domain.ID, user.ID, appID, models.Bearer, models.PASSWORD)
		if err != nil {
			return nil, err
		}
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
func (s *Store) issueTokenByClient(serviceID string, scope string) (*models.Token, error) {
	t := new(models.Token)
	t.Scope = scope
	t.CreatedAt = time.Now().Unix()
	t.ExpiresIn = s.conf.Token.ExpiresIn
	t.GrantType = models.CLIENT
	t.ServiceID = serviceID

	switch t.TokenType {
	case "bearer", "":
		t.TokenType = models.Bearer
		t.AccessToken = makeBearerToken(24)
		t.RefreshToken = makeBearerToken(32)
	case "jwt":
		t.TokenType = models.JWT
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
func (s *Store) issueTokenByRefresh(clientID, clientSecret, refreshToken string) (*models.Token, error) {
	if refreshToken == "" {
		return nil, exception.NewBadRequest("resfresh_token missed")
	}

	// 通过refresh_token找到用户的token
	old, err := s.dao.Token.GetTokenByRefresh(refreshToken)
	if err != nil {
		return nil, err
	}

	vreq := &ValidateTokenReq{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		AccessToken:  old.AccessToken,
	}

	old, err = s.ValidateTokenWithClient(vreq)
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
	tk, err := s.generateToken(old.Scope, old.DomainID, old.UserID, old.ApplicationID, models.Bearer, models.REFRESH)
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

func (s *Store) issueTokenByUpScope(clientID, clientSecret, accessToken, scope string) (*models.Token, error) {
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
	vreq := &ValidateTokenReq{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		AccessToken:  accessToken,
		Scope:        scope,
	}

	t, err := s.ValidateTokenWithClient(vreq)
	if err != nil {
		return nil, err
	}

	if t.GetProjectID() == projectID {
		return t, nil
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
	if projectID != "" {
		projectOK := false
		for i := range t.AvaliableProjects {
			if t.AvaliableProjects[i].ID == projectID {
				projectOK = true
				break
			}
		}

		if !projectOK {
			return nil, exception.NewForbidden("the project: %s not belong user: %s", projectID, t.UserID)
		}
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

func (s *Store) generateToken(scope, domainID, userID, appID string, tp models.TokenType, gt models.GrantType) (*models.Token, error) {
	t := &models.Token{
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
	case models.Bearer:
		t.AccessToken = makeBearerToken(24)
		t.RefreshToken = makeBearerToken(32)
	case models.JWT:
	default:
		return nil, exception.NewInternalServerError("unknown token type, %s, only support bearer", tp)
	}

	if err := s.dao.Token.SaveToken(t); err != nil {
		return nil, err
	}

	return t, nil
}

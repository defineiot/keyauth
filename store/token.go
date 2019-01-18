package store

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/defineiot/keyauth/dao/token"
	"github.com/defineiot/keyauth/dao/user"
	"github.com/defineiot/keyauth/internal/exception"
)

const (
	// CODE define the type of authorization request
	CODE ResponseType = "code"
	// TOKEN define the type of authorization request
	TOKEN ResponseType = "token"

	// PKCE_PLAIN is oauth pkce extension
	PKCE_PLAIN = "plain"
	// PKCE_S256 is oauth pkce extension
	PKCE_S256 = "S256"
)

// ResponseType the type of authorization request
type ResponseType string

var (
	pkceMatcher = regexp.MustCompile("^[a-zA-Z0-9~._-]{43,128}$")
)

// AuthRequest Authorize request information
type AuthRequest struct {
}

// TokenRequest use to request access token
type TokenRequest struct {
	Scope               string
	GrantType           token.GrantType
	ClientID            string
	ClientSecret        string
	AuthorizationHeader string
	DomainName          string
	Username            string
	Password            string
	Code                string
	RedirectURI         string
	RefreshToken        string
	AccessToken         string
}

// Validate validate the request
func (t *TokenRequest) Validate() error {
	switch t.GrantType {
	case token.AUTHCODE:
		if t.Code == "" {
			return exception.NewBadRequest("if use %s grant type, code is needed", t.GrantType)
		}
		if t.RedirectURI == "" {
			return exception.NewBadRequest("if use %s grant type, redirect uri is need", t.GrantType)
		}
		goto CHECK_CLIENT

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
		goto CHECK_CLIENT

	case token.CLIENT:
		goto CHECK_CLIENT

	case token.REFRESH:
		if t.RefreshToken == "" {
			return exception.NewBadRequest("if use %s grant type, refresh token is needed", t.GrantType)
		}
		goto CHECK_CLIENT

	case token.UPSCOPE:
		if t.AccessToken == "" {
			return exception.NewBadRequest("if use %s grant type, access_token is needed", t.GrantType)
		}
		if t.Scope == "" {
			return exception.NewBadRequest("if use %s grant type, scope project_id is needed", t.GrantType)
		}
		goto CHECK_CLIENT

	default:
		return exception.NewBadRequest("invalid_grant, supported grant type: [authorization_code, implicit, password, client_credentials, refresh_token, upgrade_scope]")
	}

CHECK_CLIENT:
	if t.ClientID == "" || t.ClientSecret == "" {
		return exception.NewBadRequest("if use %s grant type, client id and client secret is needed", t.GrantType)
	}

	return nil
}

// IssueToken use to issue access token
func (s *Store) IssueToken(req *TokenRequest) (*token.Token, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	app, err := s.dao.Application.GetApplicationByClientID(req.ClientID)
	if err != nil {
		return nil, err
	}

	if req.GrantType != token.IMPLICIT {
		if req.ClientSecret != app.ClientSecret {
			return nil, exception.NewUnauthorized("unauthorized_client")
		}
	}

	var t *token.Token
	switch req.GrantType {
	case token.AUTHCODE:
		t, err = s.issueTokenByAuthCode(req.ClientID, req.Code, req.RedirectURI)
		goto DEAL_ERROR
	case token.IMPLICIT:
		t, err = s.issueTokenByImplicit(req.ClientID, req.RedirectURI)
		goto DEAL_ERROR
	case token.PASSWORD:
		t, err = s.issueTokenByPassword(req.Scope, app.ClientID, req.Username, req.Password)
		goto DEAL_ERROR
	case token.CLIENT:
		t, err = s.issueTokenByClient(req.ClientID, req.Scope)
		goto DEAL_ERROR
	case token.REFRESH:
		t, err = s.issueTokenByRefresh(req.RefreshToken)
		goto DEAL_ERROR
	case token.UPSCOPE:
		t, err = s.issueTokenByUpScope(req.AccessToken, req.Scope, req.Scope)
	default:
		return nil, exception.NewBadRequest("invalid_grant only support [authorization_code, implicit, password, client_credentials, refresh_token, upgrade_scope]")
	}

DEAL_ERROR:
	if err != nil {
		return nil, err
	}

	return t, nil
}

// ValidateToken use to validate token
func (s *Store) ValidateToken(accessToken string) (*token.Token, error) {
	var err error

	tk := new(token.Token)

	// 1. get token from cache, and validate is expire
	cacheKey := "token_" + accessToken
	if s.isCache {
		if s.cache.Get(cacheKey, tk) {
			s.log.Debug("get token from cache key: %s", cacheKey)
			ok, exp := tk.IsExpired()
			if !ok {
				return nil, exception.NewExpired("token has expired, access_token: %s", tk.AccessToken)
			}
			tk.ExpiresIn = exp
			return tk, nil
		}
		s.log.Debug("get token from cache failed, key: %s", cacheKey)
	}

	// 2. get token from backend, and validate is expire
	tk, err = s.dao.Token.GetToken(accessToken)
	if err != nil {
		return nil, err
	}
	ok, delta := tk.IsExpired()
	if !ok {
		return nil, exception.NewExpired("token has expired, access_token: %s", tk.AccessToken)
	}
	tk.ExpiresIn = delta

	// 3. if this token is for user, add valiable projects
	// if tk.UserID != "" {
	// 	if tk.Scope == "" {
	// 		tk.Scope = new(token.Scope)
	// 	}
	// 	projectIDs, err := s.user.ListUserProjects(tk.DomainID, tk.UserID)
	// 	if err != nil {
	// 		return nil, exception.NewInternalServerError(err.Error())
	// 	}
	// 	tk.Scope.AvaliableProjects = projectIDs

	// 	// 4. if user hasn't work project set his default
	// 	user, err := s.user.GetUserByID(tk.UserID)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if tk.Scope.WorkProject == "" && user.DefaultProjectID != "" {
	// 		tk.Scope.WorkProject = user.DefaultProjectID
	// 	}

	// 	// 5. add user's roles
	// 	roles, err := s.user.ListUserRoles(tk.DomainID, tk.UserID)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	tk.Roles = roles
	// }

	// 4. set to cache
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

	cacheKey := "token_" + accessToken

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
func (s *Store) issueTokenByPassword(scope, clientID, account, password string) (*token.Token, error) {
	// 查询用户是否存在
	user, err := s.dao.User.GetUser(user.Account, account)
	if err != nil {
		return nil, err
	}

	// 1. vilidate user pass
	if s.hmacHash(password) != user.Password.Password {
		return nil, exception.NewForbidden("username or password invalidate")
	}

	projectIDs, err := s.dao.Project.ListUserProjects(user.Domain.ID, user.ID)
	if err != nil {
		return nil, exception.NewInternalServerError(err.Error())
	}
	fmt.Println(projectIDs)
	// if scope != nil && scope.WorkProject != "" {
	// 	if len(projectIDs) == 0 {
	// 		return nil, exception.NewForbidden("the scope project: %s, but here is no avaliable project for your", scope.WorkProject)
	// 	}
	// 	ok, err := s.checkInProject(scope.WorkProject, projectIDs)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if !ok {
	// 		return nil, exception.NewBadRequest("the project: %s not belong user: %s", scope.WorkProject, uid)
	// 	}
	// }

	// 4. generate tokend
	tk, err := s.generateToken(scope, user.Domain.ID, user.ID, clientID)
	if err != nil {
		return nil, err
	}

	// 5. add valiable projects (default project and other projects)
	// if tk.Scope == nil {
	// 	tk.Scope = new(token.Scope)
	// }
	// tk.Scope.AvaliableProjects = projectIDs

	// 6. if user hasn't work project set his default
	// if tk.Scope.WorkProject == "" && user.DefaultProjectID != "" {
	// 	tk.Scope.WorkProject = user.DefaultProjectID
	// }

	// 7. update user roles
	// roles, err := s.user.ListUserRoles(user.DomainID, user.ID)
	// if err != nil {
	// 	return nil, err
	// }
	// tk.Roles = roles

	return tk, nil
}

// issueTokenByClient implement Client Credentials Grant
// https://tools.ietf.org/html/rfc6749#section-4.4.2
func (s *Store) issueTokenByClient(clientID string, scope string) (*token.Token, error) {
	svr, err := s.dao.Service.GetServiceByClientID(clientID)
	fmt.Println(svr)
	if err != nil {
		return nil, exception.NewUnauthorized("only service client can issue by client credentials grant, but %s", err)
	}

	t := new(token.Token)
	t.Scope = scope
	t.CreatedAt = time.Now().Unix()
	t.ExpiresIn = s.conf.Token.ExpiresIn
	t.GrantType = token.CLIENT
	t.ApplicationID = clientID

	switch t.TokenType {
	case "bearer":
		t.AccessToken = makeBearerToken(24)
		t.RefreshToken = makeBearerToken(32)
	case "jwt":
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

	// 1. get old token
	old, err := s.dao.Token.GetTokenByRefresh(refreshToken)
	if err != nil {
		return nil, err
	}
	if old.UserID == "" {
		return nil, exception.NewBadRequest("the token not an user token, can't refresh")
	}

	// 2. delete old token
	if err := s.dao.Token.DeleteTokenByRefresh(refreshToken); err != nil {
		return nil, err
	}

	// 3. delete old token cache
	if s.isCache {
		cacheKey := "token_" + old.AccessToken
		if !s.cache.Delete(cacheKey) {
			s.log.Debug("delete token from cache failed, key: %s", cacheKey)
		}
		s.log.Debug("delete token from cache success, key: %s", cacheKey)
	}

	// 4. generate new token
	tk, err := s.generateToken(old.Scope, old.DomainID, old.UserID, old.ApplicationID)
	if err != nil {
		return nil, err
	}

	// 5. add valiable projects
	// projectIDs, err := s.dao.Project.ListUserProjects(tk.DomainID, tk.UserID)
	// if err != nil {
	// 	return nil, exception.NewInternalServerError(err.Error())
	// }

	// 6. add valiable projects (default project and other projects)
	// if tk.Scope != nil {
	// 	tk.Scope = new(token.Scope)
	// }
	// tk.Scope.AvaliableProjects = projectIDs

	// 7. if user hasn't work project set his default
	// user, err := s.dao.User.GetUser(user.UserID, old.UserID)
	// if err != nil {
	// 	return nil, err
	// }

	// if tk.Scope.WorkProject == "" && user.DefaultProjectID != "" {
	// 	tk.Scope.WorkProject = user.DefaultProjectID
	// }

	return tk, nil
}

func (s *Store) issueTokenByUpScope(accessToken, workProject, workDomain string) (*token.Token, error) {
	if accessToken == "" {
		return nil, exception.NewBadRequest("access_token missed")
	}
	if workProject == "" && workDomain == "" {
		return nil, exception.NewBadRequest("scope project_id or domain_id missed")
	}

	// 1. validate access token
	t, err := s.ValidateToken(accessToken)
	if err != nil {
		return nil, err
	}
	// if t.Scope.WorkProject != "" {
	// 	return nil, exception.NewBadRequest("your token is an scope token, token scope: project_id: %s", t.Scope.WorkProject)
	// }

	// 2. check the work domain is user's other domain
	// if workDomain != "" {
	// 	var validateD bool
	// 	otherDs, err := s.user.ListUserOtherDomains(t.UserID)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	for _, ad := range otherDs {
	// 		if workDomain == ad {
	// 			validateD = true
	// 			break
	// 		}
	// 	}
	// 	if !validateD {
	// 		return nil, exception.NewForbidden("the domain: %s not belong user: %s", workDomain, t.UserID)
	// 	}
	// }

	// 3. check the scope, check the scope project is belong to user
	// var pids []string
	// if workProject != "" {
	// 	projects, err := s.dao.Project.ListUserProjects(t.DomainID, t.UserID)
	// 	if err != nil {
	// 		return nil, exception.NewInternalServerError(err.Error())
	// 	}
	// 	ok, err := s.checkInProject(workProject, projects)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if !ok {
	// 		return nil, exception.NewBadRequest("the project: %s not belong user: %s", workProject, t.UserID)
	// 	}

	// 	pids = projectIDs
	// }

	// 4. generate an new token
	// scope := token.Scope{WorkProject: workProject}
	newTK, err := s.generateToken("", t.DomainID, t.UserID, t.ApplicationID)
	if err != nil {
		return nil, err
	}

	// 5. change domain
	if workDomain != "" {
		newTK.DomainID = workDomain
	}

	// 6. add aviable projects
	// if newTK.Scope == nil {
	// 	newTK.Scope = new(token.Scope)
	// }
	// newTK.Scope.AvaliableProjects = pids

	return newTK, nil
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

func (s *Store) checkInProject(targetProject string, projectIDs []string) (bool, error) {
	validated := false

	for _, p := range projectIDs {
		ok, err := s.dao.Project.CheckProjectIsExistByID(p)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, exception.NewBadRequest("project %s not exist", p)
		}

		if targetProject == p {
			validated = true
			break
		}
	}

	return validated, nil
}

func (s *Store) generateToken(scope, domainID, userID, clientID string) (*token.Token, error) {
	t := new(token.Token)
	t.Scope = scope
	t.DomainID = domainID
	t.CreatedAt = time.Now().Unix()
	t.ExpiresIn = s.conf.Token.ExpiresIn
	t.GrantType = token.PASSWORD
	t.UserID = userID
	t.ApplicationID = clientID

	switch t.TokenType {
	case "bearer":
		t.AccessToken = makeBearerToken(24)
		t.RefreshToken = makeBearerToken(32)
	case "jwt":
	default:
		return nil, exception.NewInternalServerError("unknown token type, %s", t.TokenType)
	}

	if err := s.dao.Token.SaveToken(t); err != nil {
		return nil, err
	}

	// check user role, add to token
	user, err := s.dao.User.GetUser(user.UserID, userID)
	if err != nil {
		return nil, exception.NewInternalServerError(err.Error())
	}
	for _, rn := range user.RoleNames {
		switch rn {
		case "system_admin":
			t.IsSystemAdmin = true
		case "domain_admin":
			t.IsDomainAdmin = true
		}
	}

	return t, nil
}

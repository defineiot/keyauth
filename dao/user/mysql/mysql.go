package mysql

import (
	"database/sql"

	"github.com/defineiot/keyauth/dao/user"
	"github.com/defineiot/keyauth/internal/exception"
	"github.com/defineiot/keyauth/internal/log"
	"github.com/defineiot/keyauth/internal/tools"
)

const (
	SaveUserOtherDomain   = "save-user-other-domain"
	FindUserOtherDomain   = "find-user-other-domain"
	DeleteUserOtherDomain = "delete-user-other-domain"

	SaveInvitationsRecord         = "save-invitation-record"
	UpdateInvitationsRecord       = "update-invitation-record"
	DeleteInvitationRecord        = "delete-invitation-record"
	FindUserAllInvitationsRecords = "find-user-all-invitation-record"
	FindOneInvitationRecord       = "find-one-invitation-record"

	SaveUser          = "save-user"
	SavePass          = "save-pass"
	FindDomainUsers   = "find-domain-users"
	FindUserByID      = "find-user-by-id"
	FindUserByAccount = "find-user-by-account"
	FindUserPassword  = "find-user-password"
	DeleteUserByID    = "delete-user-by-id"

	FindUserIDByName       = "find-user-id-by-name"
	FindGlobalUserIDByName = "find-global-user-id-by-name"
	BindRole               = "bind-role-to-user"
	UnbindRole             = "unbind-role-from-user"
	FindUserRoles          = "find-user-roles"

	FindUserProjects       = "find-user-projects"
	SetUserDefaultProject  = "set-user-default-project"
	SetUserPassword        = "set-user-password"
	AddProjectToUser       = "add-project-to-user"
	RemoveProjectsFromUser = "remove-projects-from-user"

	CheckUserExistByName       = "check-user-exist-by-name"
	CheckUserExistByID         = "check-user-exist-by-id"
	CheckUserGlobalExistByName = "check-user-global-exist-by-name"
	CheckUserRoleIsBind        = "check-user-role-is-bind"
)

// NewUserStore use to create domain storage service
func NewUserStore(db *sql.DB, key string, log log.IOTAuthLogger) (user.Store, error) {
	unprepared := map[string]string{
		SaveInvitationsRecord: `
			INSERT INTO invitation_records (inviter, invitee_roles, access_projects, invitation_time, expire_time, code) 
			VALUES (?,?,?,?,?,?);
		`,
		UpdateInvitationsRecord: `
			UPDATE invitation_records 
			SET invitee = ?,invitee_domain = ?,accept_time = ? 
			WHERE inviter = ? 
			AND code = ?;
		`,
		FindUserAllInvitationsRecords: `
			SELECT code, inviter, invitee, invitee_domain, invitee_roles, invitation_time, accept_time, expire_time, access_projects 
			FROM invitation_records
			WHERE inviter = ?;
		`,
		FindOneInvitationRecord: `
			SELECT code, inviter, invitee, invitee_domain, invitee_roles, invitation_time, accept_time, expire_time, access_projects  
			FROM invitation_records
			WHERE inviter = ? 
			AND code = ?;
		`,
		DeleteInvitationRecord: `
			DELETE FROM invitation_records 
			WHERE code = ?;
		`,
		SaveUser: `
			INSERT INTO users (id, department, account, mobile, email, phone, address, real_name, nick_name, gender, avatar, language, city, province, locked, domain_id, create_at, expires_active_days, default_project_id) 
			VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);
		`,
		SavePass: `
			INSERT INTO passwords (password, expires_at, create_at, user_id) 
			VALUES (?,?,?,?);
		`,
		SaveUserOtherDomain: `
			INSERT INTO user_domain_mappings (user_id, domain_id, join_at) 
			VALUES (?,?,?);
		`,
		FindUserOtherDomain: `
			SELECT domain_id 
			FROM user_domain_mappings 
			WHERE user_id = ?;
		`,
		DeleteUserOtherDomain: `
			DELETE FROM user_domain_mappings 
			WHERE user_id = ? 
			AND domain_id = ?;
		`,
		FindDomainUsers: `
			SELECT u.id, u.department, u.account, u.mobile, u.email, u.phone, u.address, u.real_name, u.nick_name, u.gender, u.avatar, u.language, u.city, u.province, u.locked, u.domain_id, u.create_at, u.expires_active_days, u.default_project_id, p.password, p.expires_at, p.create_at, p.update_at  
			FROM users u
			LEFT JOIN passwords p ON p.user_id = u.id 
			WHERE u.domain_id = ?;
		`,
		FindUserByID: `
			SELECT u.id, u.department, u.account, u.mobile, u.email, u.phone, u.address, u.real_name, u.nick_name, u.gender, u.avatar, u.language, u.city, u.province, u.locked, u.domain_id, u.create_at, u.expires_active_days, u.default_project_id, p.password, p.expires_at, p.create_at, p.update_at  
			FROM users u
			LEFT JOIN passwords p ON p.user_id = u.id 
			WHERE u.id = ?;
		`,
		FindUserByAccount: `
			SELECT u.id, u.department, u.account, u.mobile, u.email, u.phone, u.address, u.real_name, u.nick_name, u.gender, u.avatar, u.language, u.city, u.province, u.locked, u.domain_id, u.create_at, u.expires_active_days, u.default_project_id, p.password, p.expires_at, p.create_at, p.update_at  
			FROM users u
			LEFT JOIN passwords p ON p.user_id = u.id 
			WHERE u.account = ?;
		`,
		FindUserProjects: `
			SELECT m.project_id 
			FROM user_project_mappings m 
			LEFT JOIN projects p 
			ON m.project_id = p.id 
			WHERE m.user_id = ? 
			AND p.domain_id = ?;
		`,
		SetUserDefaultProject: `
			UPDATE users
			SET default_project_id = ? 
			WHERE id = ?;
		`,
		AddProjectToUser: `
			INSERT INTO user_project_mappings (user_id, project_id) 
			VALUES (?,?);
		`,
		RemoveProjectsFromUser: `
			DELETE FROM user_project_mappings 
			WHERE user_id = ? 
			AND project_id = ?;
		`,
		FindUserPassword: `
			SELECT p.password, p.expires_at, p.create_at, p.update_at 
			FROM passwords p
			WHERE user_id = ?;
		`,
		FindUserIDByName: `
			SELECT u.id 
			FROM users u
			WHERE account = ? 
			AND domain_id = ?;
		`,
		FindGlobalUserIDByName: `
		    SELECT u.id 
		    FROM users u
		    WHERE account = ?;
	     `,
		CheckUserExistByName: `
			SELECT u.account 
			FROM users u
			WHERE account = ? 
			AND domain_id = ?;
		`,
		CheckUserExistByID: `
			SELECT u.id 
			FROM users u
			WHERE id = ?;
		`,
		CheckUserGlobalExistByName: `
			SELECT u.id
			FROM users u
			WHERE account = ?;
	    `,
		BindRole: `
		    INSERT INTO role_user_mappings (domain_id, user_id, role_id) 
		    VALUES (?,?,?);
		`,
		UnbindRole: `
		    DELETE FROM role_user_mappings 
			WHERE domain_id = ?
			AND  user_id = ? 
			AND role_id = ?;
		`,
		FindUserRoles: `
			SELECT role_id 
			FROM role_user_mappings 
			WHERE domain_id = ?  
			AND user_id = ?;
		`,
		CheckUserRoleIsBind: `
		    SELECT role_id
		    FROM role_user_mappings 
			WHERE domain_id = ?
			AND user_id = ?
			AND role_id = ?;
		`,
		SetUserPassword: `
			UPDATE passwords
			SET password = ?
			WHERE user_id = ?;
		`,
	}

	// prepare all statements to verify syntax
	stmts, err := tools.PrepareStmts(db, unprepared)
	if err != nil {
		return nil, exception.NewInternalServerError("prepare user store query statment error, %s", err)
	}

	s := store{
		db:         db,
		stmts:      stmts,
		unprepared: unprepared,
		key:        key,
		log:        log,
	}

	return &s, nil
}

// DomainManager is use mongodb as storage
type store struct {
	db         *sql.DB
	stmts      map[string]*sql.Stmt
	unprepared map[string]string
	key        string
	log        log.IOTAuthLogger
}

// Close closes the database, releasing any open resources.
func (s *store) Close() error {
	return s.db.Close()
}

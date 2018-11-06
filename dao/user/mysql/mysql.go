package mysql

import (
	"database/sql"

	"github.com/defineiot/keyauth/dao/user"
	"github.com/defineiot/keyauth/internal/exception"
	"github.com/defineiot/keyauth/internal/log"
	"github.com/defineiot/keyauth/internal/tools"
)

const (
	SaveVerifyCode        = "save-verify-code"
	FindVerifyCodeByMail  = "find-verify-code-by-mail"
	FindVerifyCodeByPhone = "find-verify-code-by-phone"
	DeleteVerifyCodeByID  = "delete-verify-code-by-id"
	SaveUserOtherDomain   = "save-user-other-domain"
	FindUserOtherDomain   = "find-user-other-domain"
	DeleteUserOtherDomain = "delete-user-other-domain"

	SaveInvitationsRecord         = "save-invitation-record"
	UpdateInvitationsRecord       = "update-invitation-record"
	DeleteInvitationRecord        = "delete-invitation-record"
	FindUserAllInvitationsRecords = "find-user-all-invitation-record"
	FindOneInvitationRecord       = "find-one-invitation-record"

	SaveUser               = "save-user"
	FindAllUsers           = "find-all-users"
	FindUserByID           = "find-user-by-id"
	FindUserByName         = "find-user-by-name"
	FindUserPhones         = "find-user-phones"
	FindUserEmails         = "find-user-emails"
	FindUserPassword       = "find-user-password"
	DeleteUserByID         = "delete-user-by-id"
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
		SaveVerifyCode: `
			INSERT INTO verification_code (email_address, phone_number, code, create_at, expire_at) 
			VALUES (?,?,?,?,?);
		`,
		FindVerifyCodeByMail: `
			SELECT v.id, v.email_address, v.phone_number, v.code, v.create_at, v.expire_at
			FROM verification_code v
			WHERE code = ? 
			AND email_address = ?;
		`,
		FindVerifyCodeByPhone: `
			SELECT v.id, v.email_address, v.phone_number, v.code, v.create_at, v.expire_at
			FROM verification_code v
			WHERE code = ? 
			AND phone_number = ?;
		`,
		DeleteVerifyCodeByID: `
			DELETE FROM verification_code 
			WHERE id = ?;
		`,
		SaveInvitationsRecord: `
			INSERT INTO user_invitation_records (inviter_id, invited_user_role_ids, access_project_ids, invitation_time, expire_time, invitation_code) 
			VALUES (?,?,?,?,?,?);
		`,
		UpdateInvitationsRecord: `
			UPDATE user_invitation_records 
			SET invited_user_id = ?,invited_user_domain_id = ?,accept_time = ? 
			WHERE inviter_id = ? 
			AND invitation_code = ?;
		`,
		FindUserAllInvitationsRecords: `
			SELECT id, inviter_id, invited_user_id, invited_user_domain_id, invited_user_role_ids, invitation_time, accept_time, expire_time, invitation_code, access_project_ids
			FROM user_invitation_records
			WHERE inviter_id = ?;
		`,
		FindOneInvitationRecord: `
			SELECT id, inviter_id, invited_user_id, invited_user_domain_id, invited_user_role_ids, invitation_time, accept_time, expire_time, invitation_code, access_project_ids
			FROM user_invitation_records
			WHERE inviter_id = ? 
			AND invitation_code = ?;
		`,
		DeleteInvitationRecord: `
			DELETE FROM user_invitation_records 
			WHERE id = ?;
		`,
		SaveUser: `
			INSERT INTO users (id, name, enabled, domain_id, create_at, expires_active_days) 
			VALUES (?,?,?,?,?,?);
		`,
		SaveUserOtherDomain: `
			INSERT INTO third_party_users_domains_mapping (user_id, domain_id, create_at) 
			VALUES (?,?,?);
		`,
		FindUserOtherDomain: `
			SELECT domain_id 
			FROM third_party_users_domains_mapping 
			WHERE user_id = ?;
		`,
		DeleteUserOtherDomain: `
			DELETE FROM third_party_users_domains_mapping 
			WHERE user_id = ? 
			AND domain_id = ?;
		`,
		FindAllUsers: `
			SELECT u.id, u.name, u.enabled, u.last_active_time, u.create_at, u.expires_active_days, u.default_project_id, u.domain_id 
			FROM users u
			WHERE domain_id = ?;
		`,
		FindUserByID: `
			SELECT u.id, u.name, u.enabled, u.last_active_time, u.domain_id, u.create_at, u.expires_active_days, u.default_project_id 
			FROM users u
			WHERE id = ?;
		`,
		FindUserByName: `
			SELECT u.id, u.name, u.enabled, u.last_active_time, u.domain_id, u.create_at, u.expires_active_days, u.default_project_id 
			FROM users u 
			WHERE name = ? 
			AND domain_id = ?;
		`,
		FindUserProjects: `
			SELECT m.project_id 
			FROM users_projects_mapping m 
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
			INSERT INTO users_projects_mapping (user_id, project_id) 
			VALUES (?,?);
		`,
		RemoveProjectsFromUser: `
			DELETE FROM users_projects_mapping 
			WHERE user_id = ? 
			AND project_id = ?;
		`,
		FindUserPhones: `
			SELECT p.id, p.numbers, 'p.primary', p.description 
			FROM phones p 
			WHERE user_id = ?;
		`,
		FindUserEmails: `
			SELECT e.id, e.address, 'e.primary', e.description 
			FROM emails e
			WHERE user_id = ?;
		`,
		FindUserPassword: `
			SELECT p.password, p.expires_at, p.create_at, p.update_at 
			FROM passwords p
			WHERE user_id = ?;
		`,
		FindUserIDByName: `
			SELECT u.id 
			FROM users u
			WHERE name = ? 
			AND domain_id = ?;
		`,
		FindGlobalUserIDByName: `
		    SELECT u.id 
		    FROM users u
		    WHERE name = ?;
	     `,
		CheckUserExistByName: `
			SELECT u.name 
			FROM users u
			WHERE name = ? 
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
			WHERE name = ?;
	    `,
		BindRole: `
		    INSERT INTO roles_users_mapping (domain_id, user_id, role_name) 
		    VALUES (?,?,?);
		`,
		UnbindRole: `
		    DELETE FROM roles_users_mapping 
			WHERE domain_id = ?
			AND  user_id = ? 
			AND role_name = ?;
		`,
		FindUserRoles: `
			SELECT role_name 
			FROM roles_users_mapping 
			WHERE domain_id = ?  
			AND user_id = ?;
		`,
		CheckUserRoleIsBind: `
		    SELECT role_name 
		    FROM roles_users_mapping 
			WHERE domain_id = ?
			AND user_id = ?
			AND role_name = ?;
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

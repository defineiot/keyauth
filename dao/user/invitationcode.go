package user

// InvitationCode code
type InvitationCode struct {
	ID                   int64    `json:"-"`
	InviterID            string   `json:"inviter_id"`
	InvitedUserID        string   `json:"invited_user_id,omitempty"`
	InvitedUserDomainID  string   `json:"invited_user_domain_id,omitempty"`
	InvitedTime          int64    `json:"invited_time"`
	AcceptTime           int64    `json:"accept_time,omitempty"`
	ExpireTime           int64    `json:"expire_time,omitempty"`
	Code                 string   `json:"code"`
	InvitationURL        string   `json:"invitation_url"`
	InvitedUserRoleNames []string `json:"invited_user_role_names"`
	AccessProjects       []string `json:"access_project_ids"`
}

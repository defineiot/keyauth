package mysql

import (
	"database/sql"
	"encoding/hex"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/defineiot/keyauth/dao/user"
	"github.com/defineiot/keyauth/internal/exception"
)

func (s *store) SaveInvitationsRecord(inviterID string, invitedRoles, accessProjects []string) (*user.Invitation, error) {
	after, err := time.ParseDuration("1h")
	if err != nil {
		return nil, err
	}
	expire := time.Now().Add(after)

	code := hex.EncodeToString(uuid.NewV4().Bytes())

	ir := &user.Invitation{Code: code, Inviter: inviterID, InviteeRoles: invitedRoles, AccessProjects: accessProjects, InvitedTime: time.Now().Unix(), ExpireTime: expire.Unix()}
	_, err = s.stmts[SaveInvitationsRecord].Exec(ir.Inviter, strings.Join(invitedRoles, ","), strings.Join(accessProjects, ","), ir.InvitedTime, ir.ExpireTime, code)
	if err != nil {
		return nil, exception.NewInternalServerError("insert verify code exec sql err, %s", err)
	}

	return ir, nil
}

func (s *store) ListInvitationRecord(inviterID string) ([]*user.Invitation, error) {
	ok, err := s.CheckUserIsExistByID(inviterID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, exception.NewBadRequest("inviter user %s not exist", inviterID)
	}

	rows, err := s.stmts[FindUserAllInvitationsRecords].Query(inviterID)
	if err != nil {
		return nil, exception.NewInternalServerError("query user's invitation records error, %s", err)
	}
	defer rows.Close()

	irs := []*user.Invitation{}
	for rows.Next() {
		ir := new(user.Invitation)
		if err := rows.Scan(ir.Inviter, ir.Invitee, ir.InviteeDomain, ir.InviteeRoles, ir.InvitedTime, ir.AcceptTime, ir.ExpireTime, ir.Code, ir.AccessProjects); err != nil {
			return nil, exception.NewInternalServerError("scan user's project id error, %s", err)
		}
		irs = append(irs, ir)
	}
	return irs, nil
}

func (s *store) GetInvitationRecord(inviterID, code string) (*user.Invitation, error) {
	s.log.Debug("Get Invitation Record SQL: %s Params: inviter: %s, code: %s", s.unprepared[FindOneInvitationRecord], inviterID, code)
	ir := new(user.Invitation)
	roles := ""
	projects := ""
	err := s.stmts[FindOneInvitationRecord].QueryRow(inviterID, code).Scan(&ir.Inviter, &ir.Invitee, &ir.InviteeDomain, &roles, &ir.InvitedTime, &ir.AcceptTime, &ir.ExpireTime, &ir.Code, &projects)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("registry code %s not find", code)
		}

		return nil, exception.NewInternalServerError("query single registry code error, %s", err)
	}

	if roles != "" {
		ir.InviteeRoles = strings.Split(roles, ",")
	}
	if projects != "" {
		ir.AccessProjects = strings.Split(projects, ",")
	}

	return ir, nil
}

func (s *store) DeleteInvitationRecord(id string) error {
	ret, err := s.stmts[DeleteInvitationRecord].Exec(id)
	if err != nil {
		return exception.NewInternalServerError("delete invitation record exec sql error, %s", err)
	}
	count, err := ret.RowsAffected()
	if err != nil {
		return exception.NewInternalServerError("get delete row affected error, %s", err)
	}
	if count == 0 {
		return exception.NewBadRequest("invitation recode %s not exist", id)
	}

	return nil
}

func (s *store) UpdateInvitationsRecord(ir *user.Invitation) error {
	ir.AcceptTime = time.Now().Unix()
	_, err := s.stmts[UpdateInvitationsRecord].Exec(ir.Invitee, ir.InviteeDomain, ir.AcceptTime, ir.Inviter, ir.Code)
	if err != nil {
		return exception.NewInternalServerError("insert verify code exec sql err, %s", err)
	}

	return nil
}

package mysql

import (
	"database/sql"

	"github.com/defineiot/keyauth/dao"
	"github.com/defineiot/keyauth/dao/verifycode"
	"github.com/defineiot/keyauth/internal/exception"
	"github.com/defineiot/keyauth/internal/tools"
)

const (
	CreateVerifyCode = "create-verify-code"
	DeleteVerifyCode = "delete-verify-code"
	FindVerifyCode   = "find-verify-code"
)

// NewVerifyCodeStore use to create domain storage service
func NewVerifyCodeStore(opt *dao.Options) (verifycode.Store, error) {
	unprepared := map[string]string{
		CreateVerifyCode: `
		    INSERT INTO verification_code (code, purpose, sending_mode, sending_target, create_at, expire_at)
		    VALUES (?,?,?,?,?,?)
		`,
		FindVerifyCode: `
		    SELECT code, create_at, expire_at 
		    FROM verification_code
			WHERE purpose = ? 
			AND sending_target = ?;
	    `,
		DeleteVerifyCode: `
		    DELETE FROM verification_code
			WHERE purpose = ? 
			AND sending_target = ? 
			AND code = ?;
		`,
	}

	// prepare all statements to verify syntax
	stmts, err := tools.PrepareStmts(opt.DB, unprepared)
	if err != nil {
		return nil, exception.NewInternalServerError("prepare verify code store query statment error, %s", err)
	}

	s := store{
		db:    opt.DB,
		stmts: stmts,
	}

	return &s, nil
}

// DomainManager is use mongodb as storage
type store struct {
	db    *sql.DB
	stmts map[string]*sql.Stmt
}

// Close closes the database, releasing any open resources.
func (s *store) Close() error {
	return s.db.Close()
}

func init() {
	dao.Registe(NewVerifyCodeStore)
}

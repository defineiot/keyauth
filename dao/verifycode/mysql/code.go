package mysql

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/defineiot/keyauth/dao/verifycode"
	"github.com/defineiot/keyauth/internal/exception"
)

func (s *store) CreateVerifyCode(v *verifycode.VerifyCode) error {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	v.Code = fmt.Sprintf("%06v", rnd.Int31n(1000000))
	v.CreateAt = time.Now().Unix()

	_, err := s.stmts[CreateVerifyCode].Exec(v.Code, int(v.Purpose), int(v.SendMode), v.SendTarget, v.CreateAt, v.ExpireAt)
	if err != nil {
		return exception.NewInternalServerError("insert verify code exec sql err, %s", err)
	}

	return nil
}

func (s *store) GetVerifyCode(purpose verifycode.CodePurpose, target string) (*verifycode.VerifyCode, error) {
	v := new(verifycode.VerifyCode)
	v.Purpose = purpose
	v.SendTarget = target

	err := s.stmts[FindVerifyCode].QueryRow(int(purpose), target).Scan(&v.Code, &v.CreateAt, &v.ExpireAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("%s verify code not find", target)
		}

		return nil, exception.NewInternalServerError("query single client error, %s", err)
	}

	return v, nil
}

func (s *store) DeleteVerifyCode(purpose verifycode.CodePurpose, target string, code string) error {
	ret, err := s.stmts[DeleteVerifyCode].Exec(int(purpose), target, code)
	if err != nil {
		return exception.NewInternalServerError("delete verify code exec sql error, %s", err)
	}
	count, err := ret.RowsAffected()
	if err != nil {
		return exception.NewInternalServerError("get delete row affected error, %s", err)
	}
	if count == 0 {
		return exception.NewBadRequest("verify code %s not exist", code)
	}

	return nil
}

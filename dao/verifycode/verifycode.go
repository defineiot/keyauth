package verifycode

import (
	"github.com/defineiot/keyauth/dao/models"
)

// Store service store interface
type Store interface {
	Reader
	Writer
	Close() error
}

// Reader read service information from store
type Reader interface {
	GetVerifyCode(purpose models.CodePurpose, target string) (*models.VerifyCode, error)
}

// Writer write service information to store
type Writer interface {
	CreateVerifyCode(code *models.VerifyCode) error
	DeleteVerifyCode(purpose models.CodePurpose, target string, code string) error
}

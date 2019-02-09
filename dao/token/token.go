package token

import (
	"encoding/base64"
	"math/rand"
	"strings"
	"time"

	"github.com/defineiot/keyauth/dao/models"
)

// Store is auth service
type Store interface {
	StoreReader
	StoreWriter
	Close() error
}

// StoreReader read information from store
type StoreReader interface {
	GetUserCurrentToken(userID, appID string, gt models.GrantType) (*models.Token, error)
	GetToken(accessToken string) (*models.Token, error)
	GetTokenByRefresh(refreshToken string) (*models.Token, error)
}

// StoreWriter write information to store
type StoreWriter interface {
	SaveToken(t *models.Token) error
	DeleteTokenByRefresh(refreshToken string) error
	UpdateTokenScope(accessToken, scope string) error
	DeleteToken(accessToken string) error
}

// MakeBearerToken https://tools.ietf.org/html/rfc6750#section-2.1
// b64token    = 1*( ALPHA / DIGIT /"-" / "." / "_" / "~" / "+" / "/" ) *"="
func MakeBearerToken(lenth int) string {
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

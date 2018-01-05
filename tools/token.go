package tools

import (
	"math/rand"
	"strings"
	"time"
)

// MakeUUID use to make bearer random token
func MakeUUID(lenth int) (string, error) {
	charlist := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	password := make([]string, lenth)
	rand.Seed(time.Now().UnixNano() + int64(lenth))
	for i := 0; i < lenth; i++ {
		rn := rand.Intn(len(charlist))
		w := charlist[rn : rn+1]
		password = append(password, w)
	}

	return strings.Join(password, ""), nil
}

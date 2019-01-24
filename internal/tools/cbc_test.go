package tools_test

import (
	"testing"

	"github.com/defineiot/keyauth/internal/tools"
	"github.com/stretchr/testify/require"
)

func TestAESCBC(t *testing.T) {
	data := []byte("abcdefg")
	key := []byte("123456")

	should := require.New(t)

	cipherData, err := tools.AESCBCEncrypt(data, key)
	should.NoError(err)
	t.Logf("cipher data: %s", cipherData)

	rawData, err := tools.AESCBCDecrypt(cipherData, key)
	should.NoError(err)
	t.Logf("raw data: %s", rawData)

	should.Equal(data, []byte(rawData))
}

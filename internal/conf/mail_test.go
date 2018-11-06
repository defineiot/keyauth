package conf_test

import (
	"testing"

	"github.com/defineiot/keyauth/internal/conf"

	"github.com/stretchr/testify/assert"
)

func TestSendMail(t *testing.T) {
	mail := conf.MailConf{Host: "smtp.qq.com", Port: 465, Email: "719118794@qq.com", Password: "xxxx"}

	err := mail.Init()
	assert.NoError(t, err)

	err = mail.Send("18108053819@163.com", "喻茂峻", 1101)
	assert.NoError(t, err)
}

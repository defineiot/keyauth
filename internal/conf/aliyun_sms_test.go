package conf_test

import (
	"testing"

	"github.com/defineiot/keyauth/internal/conf"

	"github.com/stretchr/testify/assert"
)

func TestSendSMS(t *testing.T) {
	sms := &conf.AliYunSMSConf{AccessKey: "xxx", AccessSecret: "xxx", SignName: "xxx", TemplateCode: "SMS_137335150"}
	err := sms.SendSms("18108053819", "7789")
	assert.NoError(t, err)

}

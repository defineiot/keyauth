package file

import (
	"testing"

	"openauth/pkg/testutils"
	"openauth/pkg/testutils/assert"
)

const configFileContent = `
[mysql]
host = "127.0.0.1"
port = "3306"
db = "openauth"
user = "openauth"
pass = "openauth"
max_open_conn = 1000
max_idle_conn = 200
max_life_time = 60

[app]
name = "openauth"
host = "0.0.0.0"
port = "8080"
key = "this is your app key"

[log]
level = "debug"
path = "log/debug.log"

[token]
type = "bearer"
expires_in = 3600
`

// 确保文件配置为单例模式
func Test_ConfigSingleton(t *testing.T) {
	configFile := testutils.NewTempFile([]byte(configFileContent))
	defer configFile.Cleanup()
	configMgr := NewFileConf(configFile.GetPath())
	conf1, err := configMgr.GetConf()
	assert.ErrorIsNil(t, err)
	conf2, err := configMgr.GetConf()
	assert.ErrorIsNil(t, err)
	assert.Equal(t, conf1.APP.Name, conf2.APP.Name)
	conf1.APP.Name = "tempName"
	assert.Equal(t, conf1.APP.Name, conf2.APP.Name)
	t.Log("file config is singleton model.")
}

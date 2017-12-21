package env

import (
	"testing"

	"openauth/pkg/testutils"
	"openauth/pkg/testutils/assert"
)

// 确保环境配置是单例模式
func Test_EnvConfigSingleton(t *testing.T) {
	envMgr := testutils.NewEnvManager()
	defer envMgr.Cleanup()
	envMgr.Set("OA_APP_HOST", "0.0.0.0")
	envMgr.Set("OA_APP_PORT", "1234")
	envMgr.Set("OA_APP_KEY", "abc")
	envMgr.Set("OA_APP_NAME", "kk")

	envMgr.Set("OA_MYSQL_HOST", "0.0.0.0")
	envMgr.Set("OA_MYSQL_PORT", "3306")
	envMgr.Set("OA_MYSQL_USER", "user")
	envMgr.Set("OA_MYSQL_PASS", "password")
	envMgr.Set("OA_MYSQL_DB", "openauth")
	envMgr.Set("OA_MYSQL_MAX_OPEN_CONN", "1000")
	envMgr.Set("OA_MYSQL_MAX_IDEL_CONN", "1000")
	envMgr.Set("OA_MYSQL_MAX_LIFE_TIME", "1000")

	envMgr.Set("OA_LOG_FILE_PATH", "xy/abc")
	envMgr.Set("OA_LOG_LEVEL", "info")

	configMgr := NewConfigManager()
	conf1, err := configMgr.GetConf()
	assert.ErrorIsNil(t, err)
	conf2, err := configMgr.GetConf()
	assert.ErrorIsNil(t, err)
	assert.Equal(t, conf1.APP.Name, conf2.APP.Name)
	conf1.Log.FilePath = "abc/acd"
	assert.Equal(t, conf1.Log.FilePath, conf2.Log.FilePath)
	t.Log("env config is singleton model.")
}

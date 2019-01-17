package all

import (
	_ "github.com/defineiot/keyauth/dao/application/mysql"
	_ "github.com/defineiot/keyauth/dao/department/mysql"
	_ "github.com/defineiot/keyauth/dao/domain/mysql"
	_ "github.com/defineiot/keyauth/dao/project/mysql"
	_ "github.com/defineiot/keyauth/dao/role/mysql"
	_ "github.com/defineiot/keyauth/dao/service/mysql"
	_ "github.com/defineiot/keyauth/dao/token/mysql"
	_ "github.com/defineiot/keyauth/dao/user/mysql"
	_ "github.com/defineiot/keyauth/dao/verifycode/mysql"
)

# Keyauth

keyauth是一个分布式或者微服务场景下的鉴权中心, 遵循OAuth2.0规范, 参考[openstack keystone](https://developer.openstack.org/api-ref/identity/v3/?expanded=password-authentication-with-unscoped-authorization-detail,password-authentication-with-scoped-authorization-detail)和[cloud foundry uaa](http://docs.cloudfoundry.org/api/uaa/#user-token-grant-21336)设计而成, 提供如下功能:

+ 支持多租户用户管理
+ 支持OAuth2.0的中心化的身份管理
+ 支持RBAC的鉴权管理
+ 支持服务目录

具体请参考设计文档[iot-auth概要设计](./docs/design/summary.md)

## 快速开发

开发环境:

+ Golang 1.11+
+ 编辑器: 推荐使用vscod
+ 开发环境: macOS/Linux

A. 初始化数据库: sql脚本位于: cmd/ddl/schema_v1.sql

```sh
mysql -uxxxx -p < cmd/ddl/schema_v1.sql
```

B. 确认配置文件: 配置文件位于: cmd/etc/keyauth.conf, 提前配置好本地的数据库等相关配置

``` bash
[mysql]
host = "127.0.0.1"
port = "3306"
db = "keyauth"
user = "root"
pass = "passwd"
max_open_conn = 1000
max_idle_conn = 200
max_life_time = 60
```

C. 初始化系统管理员信息(仅需执行一次)

```bash
➜  keyauth git:(master) ✗ make init_admin
[INIT] 开始初始化 系统需要的角色 ...
[INIT] 创建系统管理员角色成功: system_admin
[INIT] 创建租户管理员成功: domain_admin
[INIT] 创建普通成员角色成功: member
[INIT] 开始初始化 系统管理员账户 ...
[INIT] 创建系统管理员部门成功: admin_department
[INIT] 创建系统管理员默认部门成功: default_department
[INIT] 创建系统管理员成功: admin
[INIT] 绑定系统管理员角色成功
[INIT] 绑定租户管理员角色成功
[INIT] 开始初始化 系统管理员应用 ...
[INIT] 创建Web端应用应用成功: client_id -> C1ZRpSzHM6KlhCHiC4kkML66, client_secret -> xqDAZ7kvhNGmLWppQmNZZ1vsYKtH5Nix
[INIT] 创建安卓端应用应用成功: client_id -> TR4pg4Z4FGNTDUsozMtL4f8S, client_secret -> gAtl5xfAES9ezVFW2TeHuACimeTob56s
[INIT] 创建IOS端应用应用成功: client_id -> r1n8Cjvxqy3dFakuhf5haOXL, client_secret -> usxy3jeRDZUG07V8FJE67b7DbbwcIYUl
[INIT] 创建SDK端应用应用成功: client_id -> aTgFOL7Yesq0NSqfILpKQ6A0, client_secret -> 8iZZ78kEaGag5KEf7l2UjpbsiWIJS3Y0
[INIT] 系统管理员初始化完成
```

D. 启动服务

```sh
➜  keyauth git:(master) ✗ make run
DEBU[0000] initial global variables success
DEBU[0000] registry github.com/defineiot/keyauth service features success
INFO[0000] starting keyauth service at 127.0.0.1:8080
```
# iot-auth
iot平台身份(identity)服务, 提供用户认证,鉴权,服务发现等功能。

功能:
+ 支持多租户
+ 支持OAuth2.0的中心化的身份管理
+ 支持RBAC的鉴权管理

iot-auth的具体设计请参考: [iot-auth概要设计](./docs/design/summary.md)


## 快速开发
环境要求:
+ Golang 1.10
+ 编辑器: 推荐使用vscode
+ 开发环境: macOS/Linux

根据实际情况配置好实际数据库influxDB和状态数据库redis, 配置文件位置: .keyauth/keyauth.conf
然后执行编译并启动:

```bash
$ make run
```

## 快速部署
在部署前请确认数据存储(MySQL)服务可用, 如果使用外部缓存,请确认Redis是否可用

### 二进制部署
步骤:
+ 编译二进制包
+ 启动服务

#### 编译
1. 本地编译二进制包(依赖本地已经安装好了golang环境, 并且项目在GOPATH/src目录下面)
```bash
$ make linux_build
```


2. 在容器内编译
依赖Docker运行环境, 并且已经提前下载Golang Build环境的镜像(注意配置好镜像加速源,不然拉取镜像会很慢, 推荐使用阿里的docker镜像加速):
```bash
$ docker pull golang:1.10.1
$ make docker_build
```


#### 启动
1. 将编译好的二进制包和配置文件,数据库DDL, copy到服务器上进行启动:
```bash
$ scp keyauth .keyauth/keyauth.conf ./keyauth/ddl/schema_v1.sql root@172.168.1.240:~
```

2. 初始化数据库:
```bash
$ 先在MySQL上准备好用户和库
$ mysql -h 127.0.0.1 -u iot-auth -p -D iot_auth < schema_v1.sql
```

2. 修改配置文件, 测试引导启动服务, 准备部署需要的目录和文件:
```bash
$ ./keyauth service start -f .keyauth/keyauth.conf
$ mv ./keyauth  /usr/local/bin
$ mkdir -pv /etc/keyauth
$ mv .keyauth/keyauth.conf /etc/keyauth
$ mkdir -pv /var/log/keyauth
```

3. 编写systemd的服务启动文件: /usr/lib/systemd/system/keyauth.service
```
[Unit]
Description=IOT Identity Service
After=network.target
After=network-online.target
Wants=network-online.target
Documentation=http://172.168.1.240:10080/xiniu-cloud/iot-auth
[Service]
Type=simple
WorkingDirectory=/var/log/keyauth
ExecStart=/usr/local/bin/keyauth service start -f /etc/keyauth/keyauth.conf
Restart=on-failure
RestartSec=5
LimitNOFILE=65536
[Install]
WantedBy=multi-user.target
```

4. 通过systemd启动data-gateway服务:
```sh
$ systemctl daemon-reload
$ systemctl start keyauth
$ systemctl status keyauth 
$ systemctl enable keyauth
$ systemctl is-enabled keyauth
```


3. 查看服务版本:
```bash
$ keyauth -v
Version   : v0.0.2
Build Time: 2018-05-18 21:24:45
Git Branch: master
Git Commit: fcb2dd60b6346fb0bd8944bd70514e4b59b4af56
Go Version: go1.10.1 darwin/amd64
```

### Docker部署
待更新



## 开发者手册

+ [API 文档]()
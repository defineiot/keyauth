# keyauth
[![Build Status](https://travis-ci.org/defineiot/keyauth.svg?branch=master)](https://travis-ci.org/defineiot/keyauth)
[![Go Report Card](https://goreportcard.com/badge/github.com/defineiot/keyauth)](https://goreportcard.com/report/github.com/defineiot/keyauth)

user account and authentication server with oauth 2.0


## Futures
the detail of keyauth design is here: [keyauth design summary](./docs/design.md)
+ multi tenant support
+ is an OAuth2 server that can be used for centralized identify management
+ acl for fuctions

## Usage
first you must initial dababase
```bash
$./keyauth database init
Loading config from env failed, KA_Production is not true
Loading config from file success, config: conf/keys.json
initial dabase successful
```
second start the service
```bash
$./keyauth service start
Loading config from env failed, KA_Production is not true
Loading config from file success, config: conf/keys.json
[INFO] - [Check DB Initial] - the database sql version is :1, description:初始版本
2017/08/17 17:21:04 starting keyauth service at 0.0.0.0:50000
```
# openauth
[![Build Status](https://travis-ci.org/defineiot/openauth.svg?branch=master)](https://travis-ci.org/defineiot/openauth)
[![Go Report Card](https://goreportcard.com/badge/github.com/defineiot/openauth)](https://goreportcard.com/report/github.com/defineiot/openauth)

user account and authentication server with oauth 2.0


## Futures
the detail of openauth design is here: [openauth design summary](./docs/design.md)
+ multi tenant support
+ is an OAuth2 server that can be used for centralized identify management
+ acl for fuctions

## Usage
1. first you must initial dababase
```bash
maojun@maojun-mbp  ~/GoWorkDir/src/openauth $ go build -o openauth  cmd/openauth/main.go
maojun@maojun-mbp ~/GoWorkDir/src/openauth $ ./openauth database init
initial database successful
```
2. second start the service
```bash
maojun@maojun-mbp  ~/GoWorkDir/src/openauth $ ./openauth service start
DEBU[0000] the database version: 1, desc: 初始版本           source="cmd/service.go:55"
INFO[0000] loading http middleware success               source="http/server.go:59"
INFO[0000] loading router success                        source="http/server.go:63"
INFO[0000] starting openauth service at 0.0.0.0:8080     source="http/server.go:75"
```
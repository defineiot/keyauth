# openauth

[![License][License-Image]][License-Url] [![Build][Build-Status-Image]][Build-Status-Url] [![ReportCard][ReportCard-Image]][ReportCard-Url] [![Coverage Status][Coverage-Image]][Coverage-Url] [![chat][Gitter-Image]][Gitter-Url]



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



[License-Url]: https://opensource.org/licenses/Apache-2.0
[License-Image]: https://img.shields.io/badge/license-apache2-blue.svg 
[Build-Status-Url]:https://travis-ci.org/defineiot/openauth
[Build-Status-Image]:https://travis-ci.org/defineiot/openauth.svg?branch=master
[ReportCard-Url]:https://goreportcard.com/report/github.com/defineiot/openauth
[ReportCard-Image]:https://goreportcard.com/badge/github.com/defineiot/openauth

[Coverage-Url]: https://coveralls.io/github/defineiot/openauth?branch=master
[Coverage-Image]: https://coveralls.io/repos/github/defineiot/openauth/badge.svg?branch=master
[Gitter-Url]: https://gitter.im/defineiot/Lobby
[Gitter-Image]: https://badges.gitter.im/Join_Chat.svg 
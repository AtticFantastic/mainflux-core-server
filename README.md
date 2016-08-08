# Mainflux Core Server

[![License](https://img.shields.io/badge/license-Apache%20v2.0-blue.svg)](LICENSE)
[![Build Status](https://travis-ci.org/Mainflux/mainflux-core-server.svg?branch=master)](https://travis-ci.org/Mainflux/mainflux-core-server)
[![Go Report Card](https://goreportcard.com/badge/github.com/Mainflux/mainflux-core-server)](https://goreportcard.com/report/github.com/Mainflux/mainflux-core-server)
[![Join the chat at https://gitter.im/Mainflux/mainflux](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/Mainflux/mainflux?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Mainflux Core Microservice for Mainflux IoT Platform.

### Usage

```
git clone https://github.com/drasko/go-mainflux-core-server && cd go-mainflux-core-server
go get
go build
./mainflux-core-server
```
### Dependencies
Mainflux Core Server is connected to `NATS` on northbound interface, and to `MongoDB` and `InfluxDB` southbound.

Following diagram illustrates the architecture:
![Mainflux Arch](https://github.com/Mainflux/mainflux-doc/blob/master/mermaid/arch.png)

This is why to run Mainflux Core Server you have to have running:
- [NATS](https://github.com/nats-io/gnatsd)
- [MongoDB](https://github.com/mongodb/mongo)
- [InfluxDB](https://github.com/influxdata/influxdb)

### Documentation
Development documentation can be found on our [Mainflux GitHub Wiki](https://github.com/Mainflux/mainflux/wiki).

Swagger-generated API reference can be foud at [http://mainflux.com/apidoc](http://mainflux.com/apidoc).

### Community
#### Mailing lists
- [mainflux-dev](https://groups.google.com/forum/#!forum/mainflux-dev) - developers related. This is discussion about development of Mainflux IoT cloud itself.
- [mainflux-user](https://groups.google.com/forum/#!forum/mainflux-user) - general discussion and support. If you do not participate in development of Mainflux cloud infrastructure, this is probably what you're looking for.

#### IRC
[Mainflux Gitter](https://gitter.im/Mainflux/mainflux?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

#### Twitter
[@mainflux](https://twitter.com/mainflux)

### License
[Apache License, version 2.0](LICENSE)

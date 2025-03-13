# Go gRPC Middleware

[![GitHub Releases](https://img.shields.io/github/v/release/nhatthm/go-grpc-middleware)](https://github.com/nhatthm/go-grpc-middleware/releases/latest)
[![Build Status](https://github.com/nhatthm/go-grpc-middleware/actions/workflows/test.yaml/badge.svg)](https://github.com/nhatthm/go-grpc-middleware/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/nhatthm/go-grpc-middleware/branch/master/graph/badge.svg?token=eTdAgDE2vR)](https://codecov.io/gh/nhatthm/go-grpc-middleware)
[![Go Report Card](https://goreportcard.com/badge/github.com/nhatthm/go-grpc-middleware)](https://goreportcard.com/report/github.com/nhatthm/go-grpc-middleware)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/github.com/nhatthm/go-grpc-middleware)
[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

[gRPC Go](https://github.com/grpc/grpc-go) Middleware: interceptors, helpers, utilities.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Install](#install)
- [Interceptors](#interceptors)
    - [Ctxd Logger](#ctxd-logger)
    - [Timeout](#timeout)

## Prerequisites

- `Go >= 1.23`

[<sub><sup>[table of contents]</sup></sub>](#table-of-contents)

## Install

```bash
go get github.com/nhatthm/go-grpc-middleware
```

[<sub><sup>[table of contents]</sup></sub>](#table-of-contents)

## Interceptors

### Ctxd Logger

See [bool64/ctxd](https://github.com/bool64/ctxd)

- Server middlewares
  - `ctxd.UnaryServerInterceptor`
  - `ctxd.StreamServerInterceptor`
- Client middlewares
  - `ctxd.UnaryClientInterceptor`
  - `ctxd.StreamClientInterceptor`

[<sub><sup>[table of contents]</sup></sub>](#table-of-contents)

### Timeout

There are 4 dial options for gRPC client:

- Sleep for a duration before doing the job. <br/>
  `timeout.WithStreamClientSleepInterceptor` <br/>
  `timeout.WithUnaryClientSleepInterceptor`
- Automatically creates a new context with given duration if there is none in the current context. <br/>
  `timeout.WithStreamClientTimeoutInterceptor` <br/>
  `timeout.WithUnaryClientTimeoutInterceptor` 

## Donation

If this project help you reduce time to develop, you can give me a cup of coffee :)

[<sub><sup>[table of contents]</sup></sub>](#table-of-contents)

### Paypal donation

[![paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;or scan this

<img src="https://user-images.githubusercontent.com/1154587/113494222-ad8cb200-94e6-11eb-9ef3-eb883ada222a.png" width="147px" />

[<sub><sup>[table of contents]</sup></sub>](#table-of-contents)

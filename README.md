# Balena Go

[![PkgGoDev](https://pkg.go.dev/badge/go.einride.tech/balena)](https://pkg.go.dev/go.einride.tech/balena)
[![codecov.io](https://codecov.io/github/einride/balena-go/coverage.svg?branch=master)](https://codecov.io/github/einride/balena-go?branch=master)

<br />
<img align="left" src="logo.svg" width="180" height="180">
<br />
<br />

Balena Go is a library for accessing the Balena API

Balena API docs can be found [here](https://www.balena.io/docs/reference/api/overview/)
<br />
<br />
<br />
<br />

## Install

```sh
go get go.einride.tech/balena
```

## Usage

```go
import "go.einride.tech/balena"
```

### Authentication

An [Authentication Token](https://www.balena.io/docs/reference/api/overview/#authentication)
can be used to authenticate with the API

You can then use your token to create a new client:

```go
package main

import (
    "context"
    "go.einride.tech/balena"
)

const (
    token = "mytoken"
)

func main() {
    // We supply a nil http client to make use of http.DefaultClient
    client := balena.New(nil, token)
}
```

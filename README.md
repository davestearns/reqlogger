# Request Logger

[![Build Status](https://travis-ci.org/davestearns/reqlogger.svg?branch=master)](https://travis-ci.org/davestearns/reqlogger)
[![GoDoc](https://godoc.org/github.com/davestearns/reqlogger?status.png)](https://godoc.org/github.com/davestearns/reqlogger)

This package provides an HTTP request logging middleware handler that logs all requests using the popular and efficient [zerolog](https://github.com/rs/zerolog) package.

# Installation

```bash
go get -u github.com/davestearns/reqlogger
```

# Usage

See the `example/` subdirectory for a complete example. 

```go
package main

import (
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/davestearns/reqlogger"

	"github.com/rs/zerolog"
)

const defaultAddr = ":8080"

// RootHandler handles requests for the root resource
func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func main() {
	// get the server address
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = defaultAddr
	}

	// create the logger
	var logWriter io.Writer = os.Stdout
	logPretty, _ := strconv.ParseBool(os.Getenv("LOG_PRETTY"))
	if logPretty {
		logWriter = zerolog.NewConsoleWriter()
	}
	logger := zerolog.New(logWriter).With().Timestamp().Logger()

	// create the router
	// or use https://github.com/julienschmidt/httprouter
	// for a more efficient and feature-rich router
	mux := http.NewServeMux()
	mux.HandleFunc("/", RootHandler)

	// wrap the router with the request logger middleware
	wrappedMux := reqlogger.New(mux, logger)

	// start the server
	logger.Info().Msgf("server is listening at %s", addr)
	http.ListenAndServe(":8080", wrappedMux)
}
```

Example logging output:

```json
{"level":"info","time":"2019-12-14T18:24:48-08:00","message":"server is listening at 127.0.0.1:8080"}
{"level":"info","duration":0.008072,"status":200,"bytes":13,"headers":{"user-agent":"PostmanRuntime/7.20.1"},"time":"2019-12-14T18:24:56-08:00","message":"GET /"}
{"level":"info","duration":0.005631,"status":200,"bytes":13,"headers":{"user-agent":"PostmanRuntime/7.20.1"},"time":"2019-12-14T18:24:57-08:00","message":"GET /"}
{"level":"info","duration":0.027207,"status":200,"bytes":13,"headers":{"user-agent":"PostmanRuntime/7.20.1"},"time":"2019-12-14T18:25:00-08:00","message":"GET /"}
```

## Request Headers

By default the following request headers will be logged if they are non-zero length:

- `user-agent`
- `x-request-id`
- `x-api-key`

You can alter this list using either the `.AddHeader()` method or the `.SetHeaders()` method.

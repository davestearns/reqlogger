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

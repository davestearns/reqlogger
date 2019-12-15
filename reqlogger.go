package reqlogger

import (
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// RequestLogger is an HTTP middleware handler that logs all requests.
type RequestLogger struct {
	handler    http.Handler
	logger     zerolog.Logger
	rwcPool    sync.Pool
	reqHeaders []string
}

// NewRequestLogger constructs a new RequestLogger middleware handler that wraps the
// provided `handler` to write a log message at the end of each requesting using the
// provided `logger` instance. By default, it will include the values for the following
// request headers (if present): user-agent, x-request-id, and x-api-key.
// To adjust this list, use the AddHeader() or SetHeaders() methods.
func NewRequestLogger(handler http.Handler, logger zerolog.Logger) *RequestLogger {
	return &RequestLogger{
		handler: handler,
		logger:  logger,
		rwcPool: sync.Pool{
			New: func() interface{} {
				return &ResponseWriterCapturer{}
			},
		},
		reqHeaders: []string{"x-request-id", "x-api-key", "user-agent"},
	}
}

// AddHeader adds a request header name to the list of headers that will be included in each log message.
func (rl *RequestLogger) AddHeader(name string) {
	rl.reqHeaders = append(rl.reqHeaders, name)
}

// SetHeaders resets the list of request headers to include in each log message.
func (rl *RequestLogger) SetHeaders(names []string) {
	rl.reqHeaders = names
}

// ServeHTTP processes the request and logs the results.
func (rl *RequestLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get a response writer capturer from the pool and wrap the writer
	rwc := rl.rwcPool.Get().(*ResponseWriterCapturer)
	rwc.Wrap(w)

	// capture the start time
	start := time.Now()
	// call the wrapped handler
	rl.handler.ServeHTTP(rwc, r)
	// capture the request processing duration
	duration := time.Since(start)
	// log the request
	evt := rl.logger.Info().
		Dur("duration", duration).
		Int("status", rwc.StatusCode()).
		Int("bytes", rwc.BytesWritten())

	reqHeaders := zerolog.Dict()
	for _, name := range rl.reqHeaders {
		value := r.Header.Get(name)
		if len(value) > 0 {
			reqHeaders.Str(name, value)
		}
	}
	evt.Dict("headers", reqHeaders)

	evt.Msg(r.Method + " " + r.URL.Path)

	// put the response writer capturer back into the pool
	rl.rwcPool.Put(rwc)
}

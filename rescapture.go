package reqlogger

import "net/http"

// ResponseWriterCapturer wraps an http.ResponseWriter to capture
// the response status code and number of bytes written.
type ResponseWriterCapturer struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

// NewResponseWriterCapturer constructs a new ResponseWriterCapturer,
// wrapping the provided http.ResponseWriter and capturing pertinent information.
func NewResponseWriterCapturer(w http.ResponseWriter) *ResponseWriterCapturer {
	return &ResponseWriterCapturer{
		ResponseWriter: w,
		statusCode:     http.StatusOK, // ResponseWriter defaults to OK
	}
}

// StatusCode returns the status code written by the wrapped http.ResponseWriter.
func (rwc *ResponseWriterCapturer) StatusCode() int {
	return rwc.statusCode
}

// BytesWritten returns the number of bytes written to the wrapped http.ResponseWriter.
func (rwc *ResponseWriterCapturer) BytesWritten() int {
	return rwc.bytesWritten
}

// Wrap wraps an http.ResponseWriter. This can be used to reset the
// wrapped response writer when fetching a previously-created instance
// from an instance pool.
func (rwc *ResponseWriterCapturer) Wrap(w http.ResponseWriter) {
	rwc.ResponseWriter = w
	rwc.statusCode = http.StatusOK
	rwc.bytesWritten = 0
}

// Write overrides the ResponseWriter.Write method to capture the number of bytes written.
func (rwc *ResponseWriterCapturer) Write(buf []byte) (int, error) {
	rwc.bytesWritten += len(buf)
	return rwc.ResponseWriter.Write(buf)
}

// WriteHeader overrides the ResponseWriter.WriteHeader to capture the response status code.
func (rwc *ResponseWriterCapturer) WriteHeader(statusCode int) {
	rwc.statusCode = statusCode
	rwc.ResponseWriter.WriteHeader(statusCode)
}

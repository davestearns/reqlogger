package reqlogger

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewResponseWriterCapturer(t *testing.T) {
	w := httptest.NewRecorder()
	expected := &ResponseWriterCapturer{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		bytesWritten:   0,
	}
	actual := NewResponseWriterCapturer(w)
	assert.Equal(t, expected, actual)
}

func TestWrap(t *testing.T) {
	w := httptest.NewRecorder()
	w2 := httptest.NewRecorder()
	expected := &ResponseWriterCapturer{
		ResponseWriter: w2,
		statusCode:     http.StatusOK,
		bytesWritten:   0,
	}
	actual := &ResponseWriterCapturer{
		ResponseWriter: w,
		statusCode:     http.StatusBadRequest,
		bytesWritten:   200,
	}
	actual.Wrap(w2)
	assert.Equal(t, expected, actual)
}

func TestCapturesStatusCode(t *testing.T) {
	w := httptest.NewRecorder()
	rwc := NewResponseWriterCapturer(w)
	rwc.WriteHeader(http.StatusNotFound)
	assert.Equal(t, http.StatusNotFound, rwc.StatusCode())
}

func TestCapturesBytesWritten(t *testing.T) {
	w := httptest.NewRecorder()
	rwc := NewResponseWriterCapturer(w)
	msg := []byte("testing")
	iters := 4
	for i := 0; i < iters; i++ {
		rwc.Write(msg)
	}
	assert.Equal(t, len(msg)*iters, rwc.BytesWritten())
}

func TestReturnsHeaders(t *testing.T) {
	w := httptest.NewRecorder()
	rwc := NewResponseWriterCapturer(w)
	assert.Equal(t, w.Header(), rwc.Header())
}

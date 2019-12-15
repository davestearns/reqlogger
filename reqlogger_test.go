package reqlogger

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rs/zerolog"
)

func TestDefaultReqHeaders(t *testing.T) {
	rl := New(http.NewServeMux(), zerolog.New(ioutil.Discard))
	expected := []string{"user-agent", "x-api-key", "x-request-id"}
	assert.ElementsMatch(t, expected, rl.reqHeaders)
}

func TestAddHeader(t *testing.T) {
	rl := New(http.NewServeMux(), zerolog.New(ioutil.Discard))
	rl.AddHeader("x-foo")
	expected := []string{"user-agent", "x-api-key", "x-request-id", "x-foo"}
	assert.ElementsMatch(t, expected, rl.reqHeaders)
}

func TestSetHeaders(t *testing.T) {
	rl := New(http.NewServeMux(), zerolog.New(ioutil.Discard))
	expected := []string{"x-foo"}
	rl.SetHeaders(expected)
	assert.ElementsMatch(t, expected, rl.reqHeaders)
}

func TestLogsRequest(t *testing.T) {
	responseMsg := []byte("test response")
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write(responseMsg)
	}
	buf := bytes.NewBuffer(nil)
	r := httptest.NewRequest("GET", "/test", nil)
	r.Header.Add("user-agent", "test-user-agent")
	r.Header.Add("x-api-key", "test-api-key")
	r.Header.Add("x-request-id", "test-request-id")
	r.Header.Add("x-foo", "bar")

	rl := New(http.HandlerFunc(handler), zerolog.New(buf))
	rl.AddHeader("x-foo")
	rl.ServeHTTP(httptest.NewRecorder(), r)

	assert.True(t, buf.Len() > 0)
	evt := make(map[string]interface{})
	if err := json.Unmarshal(buf.Bytes(), &evt); err != nil {
		t.Fatalf("unmarshaling wrtten zerolog event: %v", err)
	}

	assert.Equal(t, "info", evt["level"])
	assert.Equal(t, "GET /test", evt["message"])
	assert.NotZero(t, evt["duration"])
	assert.Equal(t, float64(200), evt["status"])
	assert.Equal(t, float64(len(responseMsg)), evt["bytes"])

	headers, ok := evt["headers"].(map[string]interface{})
	assert.Equal(t, true, ok, "headers were not a map")
	assert.Equal(t, 4, len(headers))
	assert.Equal(t, "test-user-agent", headers["user-agent"])
	assert.Equal(t, "test-api-key", headers["x-api-key"])
	assert.Equal(t, "test-request-id", headers["x-request-id"])
	assert.Equal(t, "bar", headers["x-foo"])
}

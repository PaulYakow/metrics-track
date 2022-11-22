package v1

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/PaulYakow/metrics-track/internal/pkg/logger"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"github.com/PaulYakow/metrics-track/internal/usecase/repo"
	"github.com/PaulYakow/metrics-track/internal/usecase/services/hasher"
)

func init() {
	if err := os.Chdir("../../../.."); err != nil {
		panic(err)
	}
}

func testRequest(t *testing.T, ts *httptest.Server, content, method, path, body string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, bytes.NewBufferString(body))
	require.NoError(t, err)

	req.Header.Set("Content-Type", content)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}

func TestURLRoutes(t *testing.T) {
	tests := []struct {
		name    string
		content string
		method  string
		path    string
		want    int
	}{
		{
			name:    "update valid gauge",
			content: "text/plain",
			method:  "POST",
			path:    "/update/gauge/testGauge/0",
			want:    200,
		},
		{
			name:    "update invalid gauge value",
			content: "text/plain",
			method:  "POST",
			path:    "/update/gauge/testGauge/null",
			want:    400,
		},
		{
			name:    "update no gauge value",
			content: "text/plain",
			method:  "POST",
			path:    "/update/gauge/testGauge/",
			want:    400,
		},
		{
			name:    "update no gauge name",
			content: "text/plain",
			method:  "POST",
			path:    "/update/gauge/",
			want:    404,
		},
		{
			name:    "update invalid content gauge value",
			content: "text/html",
			method:  "POST",
			path:    "/update/gauge/testGauge/100",
			want:    400,
		},
		{
			name:    "read valid gauge",
			content: "text/plain",
			method:  "GET",
			path:    "/value/gauge/testGauge",
			want:    200,
		},
		{
			name:    "read no gauge name",
			content: "text/plain",
			method:  "GET",
			path:    "/value/gauge/",
			want:    404,
		},
		{
			name:    "read unknown gauge name",
			content: "text/plain",
			method:  "GET",
			path:    "/value/gauge/unknownGauge",
			want:    404,
		},
		{
			name:    "read invalid content gauge value",
			content: "text/html",
			method:  "GET",
			path:    "/value/gauge/testGauge",
			want:    400,
		},
		{
			name:    "update valid counter",
			content: "text/plain",
			method:  "POST",
			path:    "/update/counter/testCounter/10",
			want:    200,
		},
		{
			name:    "update invalid counter value",
			content: "text/plain",
			method:  "POST",
			path:    "/update/counter/testCounter/null",
			want:    400,
		},
		{
			name:    "update no counter value",
			content: "text/plain",
			method:  "POST",
			path:    "/update/counter/testCounter/",
			want:    400,
		},
		{
			name:    "update no counter name",
			content: "text/plain",
			method:  "POST",
			path:    "/update/counter/",
			want:    404,
		},
		{
			name:    "update invalid content counter value",
			content: "text/html",
			method:  "POST",
			path:    "/update/counter/testCounter/100",
			want:    400,
		},
		{
			name:    "read valid counter",
			content: "text/plain",
			method:  "GET",
			path:    "/value/counter/testCounter",
			want:    200,
		},
		{
			name:    "read no counter name",
			content: "text/plain",
			method:  "GET",
			path:    "/value/counter/",
			want:    404,
		},
		{
			name:    "read unknown counter name",
			content: "text/plain",
			method:  "GET",
			path:    "/value/counter/unknownCounter",
			want:    404,
		},
		{
			name:    "read invalid content counter value",
			content: "text/html",
			method:  "GET",
			path:    "/value/counter/testCounter",
			want:    400,
		},
		{
			name:    "update unknown metric type",
			content: "text/plain",
			method:  "POST",
			path:    "/update/unknown/test/0",
			want:    501,
		},
		{
			name:    "read unknown metric type",
			content: "text/plain",
			method:  "GET",
			path:    "/value/unknown/test",
			want:    404,
		},
		{
			name:    "read all metrics",
			content: "",
			method:  "GET",
			path:    "/",
			want:    200,
		},
		{
			name:    "ping memory repo",
			content: "",
			method:  "GET",
			path:    "/ping",
			want:    500,
		},
	}

	r := NewRouter(usecase.NewServerUC(repo.NewServerMemory(), hasher.New("")), logger.New())
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := testRequest(t, ts, tt.content, tt.method, tt.path, "")
			defer resp.Body.Close()

			assert.Equal(t, tt.want, resp.StatusCode)
		})
	}
}

func TestJSONRoutes(t *testing.T) {
	tests := []struct {
		name    string
		content string
		method  string
		path    string
		body    string
		want    int
	}{
		{
			name:    "update valid gauge",
			content: "application/json",
			method:  "POST",
			path:    "/update",
			body:    `{"id": "testGauge","type": "gauge","value": 13}`,
			want:    200,
		},
		{
			name:    "update invalid gauge value",
			content: "application/json",
			method:  "POST",
			path:    "/update",
			body:    `{"id": "testGauge","type": "gauge","value": "invalid"}`,
			want:    400,
		},
		{
			name:    "update no gauge value",
			content: "application/json",
			method:  "POST",
			path:    "/update",
			body:    `{"id": "testGauge","type": "gauge"}`,
			want:    400,
		},
		{
			name:    "update no gauge name",
			content: "application/json",
			method:  "POST",
			path:    "/update",
			body:    `{"type": "gauge","value": 13}`,
			want:    404,
		},
		{
			name:    "update invalid content gauge value",
			content: "text/html",
			method:  "POST",
			path:    "/update",
			body:    `{"id": "testGauge","type": "gauge","value": 13}`,
			want:    400,
		},
		{
			name:    "read valid gauge",
			content: "application/json",
			method:  "POST",
			path:    "/value",
			body:    `{"id": "testGauge","type": "gauge"}`,
			want:    200,
		},
		{
			name:    "read no gauge name",
			content: "application/json",
			method:  "POST",
			path:    "/value",
			body:    `{"type": "gauge"}`,
			want:    404,
		},
		{
			name:    "read invalid content gauge value",
			content: "text/html",
			method:  "POST",
			path:    "/value",
			body:    `{"id": "testGauge","type": "gauge"}`,
			want:    400,
		},
		{
			name:    "update valid counter",
			content: "application/json",
			method:  "POST",
			path:    "/update",
			body:    `{"id": "testCounter","type": "counter","delta": 13}`,
			want:    200,
		},
		{
			name:    "update invalid counter value",
			content: "application/json",
			method:  "POST",
			path:    "/update",
			body:    `{"id": "testCounter","type": "counter","delta": "invalid"}`,
			want:    400,
		},
		{
			name:    "update no counter value",
			content: "application/json",
			method:  "POST",
			path:    "/update",
			body:    `{"id": "testCounter","type": "counter"}`,
			want:    400,
		},
		{
			name:    "update no counter name",
			content: "application/json",
			method:  "POST",
			path:    "/update",
			body:    `{"type": "counter","delta": 13}`,
			want:    404,
		},
		{
			name:    "update invalid content counter value",
			content: "text/html",
			method:  "POST",
			path:    "/update",
			body:    `{"id": "testCounter","type": "counter","delta": 13}`,
			want:    400,
		},
		{
			name:    "read valid counter",
			content: "application/json",
			method:  "POST",
			path:    "/value",
			body:    `{"id": "testCounter","type": "counter"}`,
			want:    200,
		},
		{
			name:    "read no counter name",
			content: "application/json",
			method:  "POST",
			path:    "/value",
			body:    `{"type": "counter"}`,
			want:    404,
		},
		{
			name:    "read invalid content counter value",
			content: "text/html",
			method:  "POST",
			path:    "/value",
			body:    `{"id": "testCounter","type": "counter"}`,
			want:    400,
		},
		{
			name:    "update unknown metric type",
			content: "application/json",
			method:  "POST",
			path:    "/update",
			body:    `{"id": "testUnknown","type": "unknown"}`,
			want:    400,
		},
		{
			name:    "read unknown metric type",
			content: "application/json",
			method:  "POST",
			path:    "/value",
			body:    `{"id": "testUnknown","type": "unknown"}`,
			want:    404,
		},
	}

	r := NewRouter(usecase.NewServerUC(repo.NewServerMemory(), hasher.New("")), logger.New())
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := testRequest(t, ts, tt.content, tt.method, tt.path, tt.body)
			defer resp.Body.Close()

			assert.Equal(t, tt.want, resp.StatusCode)
		})
	}
}

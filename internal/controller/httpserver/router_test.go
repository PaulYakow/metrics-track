package httpserver

import (
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"github.com/PaulYakow/metrics-track/internal/usecase/repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}

func TestRouter(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
		want   int
	}{
		{
			name:   "valid_update_gauge",
			method: "POST",
			path:   "/update/gauge/testGauge/0",
			want:   200,
		},
		{
			name:   "valid_update_counter",
			method: "POST",
			path:   "/update/counter/testCounter/10",
			want:   200,
		},
		{
			name:   "invalid_value",
			method: "POST",
			path:   "/update/gauge/testGauge/null",
			want:   400,
		},
		{
			name:   "no_metric_value",
			method: "POST",
			path:   "/update/gauge/testGauge/",
			want:   400,
		},
		{
			name:   "no_metric_name",
			method: "POST",
			path:   "/update/gauge/",
			want:   404,
		},
		{
			name:   "update_without_header",
			method: "POST",
			path:   "/update/",
			want:   400,
		},
		{
			name:   "unknown_type",
			method: "POST",
			path:   "/update/unknown/testCounter/0",
			want:   501,
		},
	}

	r := NewRouter(usecase.NewServerUC(repo.NewServerMemory()))
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := testRequest(t, ts, tt.method, tt.path)
			defer resp.Body.Close()

			assert.Equal(t, tt.want, resp.StatusCode)
		})
	}
}

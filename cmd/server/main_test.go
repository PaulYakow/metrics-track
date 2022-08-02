package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_gaugeHandler(t *testing.T) {
	tests := []struct {
		name    string
		request string
		want    int
	}{
		{
			name:    "valid_update",
			request: "/update/gauge/testGauge/0",
			want:    200,
		},
		{
			name:    "invalid_value",
			request: "/update/gauge/testGauge/null",
			want:    400,
		},
		{
			name:    "no_metric_value",
			request: "/update/gauge/testGauge/",
			want:    400,
		},
		{
			name:    "no_metric_name",
			request: "/update/gauge/",
			want:    404,
		},
		{
			name:    "only_update",
			request: "/update/",
			want:    404,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()
			gaugeHandler(w, request)
			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want, result.StatusCode)
		})
	}
}

func Test_counterHandler(t *testing.T) {
	tests := []struct {
		name    string
		request string
		want    int
	}{
		{
			name:    "valid_update",
			request: "/update/counter/testCounter/0",
			want:    200,
		},
		{
			name:    "invalid_value",
			request: "/update/counter/testCounter/null",
			want:    400,
		},
		{
			name:    "no_metric_value",
			request: "/update/counter/testCounter/",
			want:    400,
		},
		{
			name:    "no_metric_name",
			request: "/update/counter/",
			want:    404,
		},
		{
			name:    "only_update",
			request: "/update/",
			want:    404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()
			counterHandler(w, request)
			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want, result.StatusCode)
		})
	}
}

func Test_defaultHandler(t *testing.T) {
	tests := []struct {
		name    string
		request string
		want    int
	}{
		{
			name:    "unknown_type",
			request: "/update/unknown/testCounter/0",
			want:    501,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()
			defaultHandler(w, request)
			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want, result.StatusCode)
		})
	}
}

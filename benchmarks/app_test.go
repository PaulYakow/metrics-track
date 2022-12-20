package benchmarks

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	v1 "github.com/PaulYakow/metrics-track/internal/controller/server/v1"
	"github.com/PaulYakow/metrics-track/internal/pkg/logger"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"github.com/PaulYakow/metrics-track/internal/usecase/repo"
	"github.com/PaulYakow/metrics-track/internal/usecase/services/hasher"
)

func init() {
	if err := os.Chdir("../"); err != nil {
		panic(err)
	}
}

func BenchmarkV1URLRoutes(b *testing.B) {
	r := v1.NewRouter(usecase.NewServerUC(repo.NewServerMemory(), hasher.New("")), logger.New(), nil)
	w := httptest.NewRecorder()

	b.ResetTimer()

	b.Run("POST gauge", func(b *testing.B) {
		b.StopTimer()
		req := httptest.NewRequest(http.MethodPost, "/update/gauge/urlGauge/0", nil)
		req.Header.Set("Content-Type", "text/plain")

		b.ReportAllocs()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			r.ServeHTTP(w, req)
		}
	})

	b.Run("POST counter", func(b *testing.B) {
		b.StopTimer()
		req := httptest.NewRequest(http.MethodPost, "/update/counter/urlCounter/1", nil)
		req.Header.Set("Content-Type", "text/plain")

		b.ReportAllocs()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			r.ServeHTTP(w, req)
		}
	})

	b.Run("GET gauge", func(b *testing.B) {
		b.StopTimer()
		req := httptest.NewRequest(http.MethodGet, "/value/gauge/urlGauge", nil)
		req.Header.Set("Content-Type", "text/plain")

		b.ReportAllocs()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			r.ServeHTTP(w, req)
		}
	})

	b.Run("GET counter", func(b *testing.B) {
		b.StopTimer()
		req := httptest.NewRequest(http.MethodGet, "/value/counter/urlCounter", nil)
		req.Header.Set("Content-Type", "text/plain")

		b.ReportAllocs()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			r.ServeHTTP(w, req)
		}
	})

	b.Run("GET ping", func(b *testing.B) {
		b.StopTimer()
		req := httptest.NewRequest(http.MethodGet, "/ping", nil)

		b.ReportAllocs()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			r.ServeHTTP(w, req)
		}
	})

	b.Run("GET all", func(b *testing.B) {
		b.StopTimer()
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		b.ReportAllocs()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			r.ServeHTTP(w, req)
		}
	})
}

func BenchmarkV1JSONRoutes(b *testing.B) {
	r := v1.NewRouter(usecase.NewServerUC(repo.NewServerMemory(), hasher.New("")), logger.New(), nil)
	w := httptest.NewRecorder()

	b.ResetTimer()

	b.Run("POST gauge", func(b *testing.B) {
		b.StopTimer()
		data := `{"id": "testGauge", "type": "gauge", "value": 13}`
		req := httptest.NewRequest(http.MethodPost, "/update", bytes.NewBufferString(data))
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")

		b.ReportAllocs()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			r.ServeHTTP(w, req)
		}
	})

	b.Run("POST counter", func(b *testing.B) {
		b.StopTimer()
		data := `{"id": "testCounter", "type": "counter", "value": 1}`
		req := httptest.NewRequest(http.MethodPost, "/update", bytes.NewBufferString(data))
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")

		b.ReportAllocs()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			r.ServeHTTP(w, req)
		}
	})

	b.Run("GET gauge", func(b *testing.B) {
		b.StopTimer()
		data := `{"id": "testGauge", "type": "gauge"}`
		req := httptest.NewRequest(http.MethodPost, "/value", bytes.NewBufferString(data))
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")

		b.ReportAllocs()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			r.ServeHTTP(w, req)
		}
	})

	b.Run("GET counter", func(b *testing.B) {
		b.StopTimer()
		data := `{"id": "testCounter", "type": "counter"}`
		req := httptest.NewRequest(http.MethodPost, "/value", bytes.NewBufferString(data))
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")

		b.ReportAllocs()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			r.ServeHTTP(w, req)
		}
	})
}

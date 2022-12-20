package v1

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func compressGzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func (s *serverRoutes) decryptData(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			s.logger.Error(fmt.Errorf("decrypt - read body %q: %w", r.URL.Path, err))
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		if s.decoder != nil {
			if body, err = s.decoder.Decrypt(body); err != nil {
				s.logger.Fatal(err)
			}
		}

		ctx := context.WithValue(r.Context(), "body", body)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}

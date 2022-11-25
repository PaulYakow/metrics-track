package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/PaulYakow/metrics-track/internal/entity"
)

func (s *serverRoutes) createUpdateRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/", s.updateByJSON)

	r.Route("/{type}", func(r chi.Router) {
		r.Post("/", func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusNotFound)
		})
		r.Post("/{name}/", func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusBadRequest)
		})
		r.Post("/{name}/{value}", s.updateByURL)
	})

	return r
}

func (s *serverRoutes) updateByURL(rw http.ResponseWriter, r *http.Request) {
	if !isContentTypeMatch(r, "text/plain") && !isContentTypeMatch(r, "") {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	mType := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")
	rawValue := chi.URLParam(r, "value")

	m, err := entity.Create(mType, name, rawValue)
	if err != nil {
		s.logger.Error(fmt.Errorf("router - create temp metric: %w", err))

		if errors.Is(err, entity.ErrParseValue) {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		if errors.Is(err, entity.ErrUnknownType) {
			rw.WriteHeader(http.StatusNotImplemented)
			return
		}

		return
	}

	if err = s.uc.Save(m); err != nil {
		s.logger.Error(fmt.Errorf("router - save metric: %w", err))
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (s *serverRoutes) updateByJSON(rw http.ResponseWriter, r *http.Request) {
	// Обработать JSON из тела запроса - сохранить в соответствующую метрику переданное значение
	if !isContentTypeMatch(r, "application/json") {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Error(fmt.Errorf("read request body %q: %w", r.URL.Path, err))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	reqMetric := entity.Metric{}
	if err = json.Unmarshal(body, &reqMetric); err != nil {
		s.logger.Error(fmt.Errorf("router - update metric: %q (%v)", err, *r))
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if reqMetric.ID == "" {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	if reqMetric.Value == nil && reqMetric.Delta == nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	s.logger.Info("router - POST /update: %v", reqMetric)

	if err = s.uc.Save(&reqMetric); err != nil {
		s.logger.Error(fmt.Errorf("router - save value to storage: %q", err))

		if errors.Is(err, entity.ErrHashMismatch) {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

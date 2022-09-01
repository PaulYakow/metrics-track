package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PaulYakow/metrics-track/internal/entity"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
)

type updates struct{}

func (rs updates) Routes(s *serverRoutes) chi.Router {
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
	mType := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")
	rawValue := chi.URLParam(r, "value")

	m, err := entity.Create(mType, name, rawValue)
	if err != nil {
		s.logger.Error(fmt.Errorf("router - create temp metric: %w", err))

		if errors.Is(err, entity.ErrParseValue) {
			rw.WriteHeader(http.StatusBadRequest)
		}

		if errors.Is(err, entity.ErrUnknownType) {
			rw.WriteHeader(http.StatusNotImplemented)
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
	if r.Header.Get("Content-Type") != "application/json" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("read request body %q: %v", r.URL.Path, err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// todo: повторяется дважды (в values) - вынести в отдельную функцию (возможно убрать в метод самой метрики)
	reqMetric := entity.Metric{}
	if err = json.Unmarshal(body, &reqMetric); err != nil {
		s.logger.Error(fmt.Errorf("router - update metric: %q", err))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.logger.Info("router - request update metric: %v", reqMetric)

	if err = s.uc.Save(&reqMetric); err != nil {
		s.logger.Error(fmt.Errorf("router - save value to storage: %q", err))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

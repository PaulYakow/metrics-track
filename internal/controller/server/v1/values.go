package v1

import (
	"encoding/json"
	"fmt"
	"github.com/PaulYakow/metrics-track/internal/entity"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

type values struct{}

func (rs values) Routes(s *serverRoutes) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/", s.valueByJSON)
	r.Get("/{type}/{name}", s.valueByURL)

	return r
}

func (s *serverRoutes) valueByURL(rw http.ResponseWriter, r *http.Request) {
	mType := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")

	metric, err := s.uc.Get(entity.Metric{
		ID:    name,
		MType: mType,
	})
	if err != nil {
		s.logger.Error(err)
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	respBody := []byte(metric.GetValue())
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(respBody)
}

func (s *serverRoutes) valueByJSON(rw http.ResponseWriter, r *http.Request) {
	// 1. В теле запроса JSON с ID и MType
	// 2. Заполнить значение метрики
	// 3. Отправить ответный JSON
	if r.Header.Get("Content-Type") != "application/json" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Error(fmt.Errorf("router - read request body %q: %v", r.URL.Path, err))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	metric := entity.Metric{}
	if err = json.Unmarshal(reqBody, &metric); err != nil {
		s.logger.Error(fmt.Errorf("router - read value by json (unmarshal): %v", err))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.logger.Info("router - request get metric: %v", metric)

	value, err := s.uc.Get(metric)
	if err != nil {
		s.logger.Error(fmt.Errorf("router - get value: %v", err))
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	respBody, err := json.Marshal(&value)
	if err != nil {
		s.logger.Error(fmt.Errorf("router - marshal to json: %v", err))
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(respBody)
}

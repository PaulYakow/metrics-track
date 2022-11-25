package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/PaulYakow/metrics-track/internal/entity"
)

func (s *serverRoutes) createValueRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/", s.valueByJSON)
	r.Get("/{type}/{name}", s.valueByURL)

	return r
}

func (s *serverRoutes) valueByURL(rw http.ResponseWriter, r *http.Request) {
	if !isContentTypeMatch(r, "text/plain") && !isContentTypeMatch(r, "") {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	mType := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")

	metric, err := s.uc.Get(r.Context(),
		entity.Metric{
			ID:    name,
			MType: mType,
		})
	if err != nil {
		s.logger.Error(err)
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	respBody := metric.GetValue()
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(respBody)
}

func (s *serverRoutes) valueByJSON(rw http.ResponseWriter, r *http.Request) {
	// 1. В теле запроса JSON с ID и MType
	// 2. Заполнить значение метрики
	// 3. Отправить ответный JSON
	if !isContentTypeMatch(r, "application/json") {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Error(fmt.Errorf("router - read request body %q: %v", r.URL.Path, err))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	reqMetric := entity.Metric{}
	if err = json.Unmarshal(body, &reqMetric); err != nil {
		s.logger.Error(fmt.Errorf("router - read value by json (unmarshal): %v", err))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if reqMetric.ID == "" {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	s.logger.Info("router - POST /value: %v", reqMetric)

	metric, err := s.uc.Get(r.Context(), reqMetric)
	if err != nil {
		s.logger.Error(fmt.Errorf("router - get value: %v", err))
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	respBody, err := json.Marshal(&metric)
	if err != nil {
		s.logger.Error(fmt.Errorf("router - marshal to json: %v", err))
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(respBody)
}

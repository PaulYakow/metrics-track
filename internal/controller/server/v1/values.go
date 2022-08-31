package v1

import (
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
)

type values struct{}

func (rs values) Routes(s *serverRoutes) chi.Router {
	r := chi.NewRouter()

	r.Post("/", s.valueByJSON)
	r.Get("/{type}/{name}", s.valueByURL)

	return r
}

func (s *serverRoutes) valueByURL(rw http.ResponseWriter, r *http.Request) {
	mType := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")

	value, err := s.uc.GetValueByType(mType, name)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	respBody := []byte(value)
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
		log.Printf("read request body %q: %v", r.URL.Path, err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	respBody, err := s.uc.GetValueByJSON(reqBody)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(respBody)
}

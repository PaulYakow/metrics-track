package v1

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type updates struct{}

func (rs updates) Routes(s *serverRoutes) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/", s.updateByJSONBatch)

	return r
}

func (s *serverRoutes) updateByJSONBatch(rw http.ResponseWriter, r *http.Request) {

}

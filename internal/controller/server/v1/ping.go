package v1

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type ping struct{}

func (rs ping) Routes(s *serverRoutes) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", s.pingDb)

	return r
}

func (s *serverRoutes) pingDb(rw http.ResponseWriter, r *http.Request) {
	err := s.uc.CheckRepo()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

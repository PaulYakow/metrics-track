package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/PaulYakow/metrics-track/internal/entity"
)

const defaultBatchCap = 20

func (s *serverRoutes) createBatchUpdateRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(s.decryptData)

	r.Post("/", s.updateByJSONBatch)

	return r
}

func (s *serverRoutes) updateByJSONBatch(rw http.ResponseWriter, r *http.Request) {
	if !isContentTypeMatch(r, "application/json") {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	body := r.Context().Value(bodyCtxKey).([]byte)

	var err error
	reqMetrics := make([]entity.Metric, 0, defaultBatchCap)
	if err = json.Unmarshal(body, &reqMetrics); err != nil {
		s.logger.Error(fmt.Errorf("router - batch update unmarshal: %w", err))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.logger.Info("router - POST /updates: %v", reqMetrics)

	if err = s.uc.SaveBatch(reqMetrics); err != nil {
		s.logger.Error(fmt.Errorf("router - batch update save to storage: %w", err))

		if errors.Is(err, entity.ErrHashMismatch) {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

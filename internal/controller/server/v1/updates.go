package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PaulYakow/metrics-track/internal/entity"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

const _defaultBatchCap = 20

type updates struct{}

func (rs updates) Routes(s *serverRoutes) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/", s.updateByJSONBatch)

	return r
}

func (s *serverRoutes) updateByJSONBatch(rw http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Error(fmt.Errorf("router - batch update read body %q: %w", r.URL.Path, err))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// todo: повторяется (в values) - вынести в отдельную функцию (возможно убрать в метод самой метрики)
	var reqMetrics = make([]entity.Metric, 0, _defaultBatchCap)
	if err = json.Unmarshal(body, &reqMetrics); err != nil {
		s.logger.Error(fmt.Errorf("router - batch update unmarshal: %w", err))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.logger.Info("router - request batch update: %v", reqMetrics)

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

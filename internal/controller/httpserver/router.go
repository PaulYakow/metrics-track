package httpserver

import (
	"compress/flate"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
)

type serverRoutes struct {
	uc usecase.IServer
}

func NewRouter(uc usecase.IServer) chi.Router {
	s := &serverRoutes{uc: uc}

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Compress(flate.BestCompression))

	mux.Route("/", func(r chi.Router) {
		r.Get("/", s.getListOfMetrics)

		r.Post("/value/", s.postValueByJSON)
		r.Get("/value/{type}/{name}", s.getMetricValue)

		r.Post("/update/", s.postUpdateByJSON)

		r.Post("/update/gauge/", func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusNotFound)
		})
		r.Post("/update/gauge/{name}/", func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusBadRequest)
		})
		r.Post("/update/gauge/{name}/{value}", s.postGauge)

		r.Post("/update/counter/", func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusNotFound)
		})
		r.Post("/update/counter/{name}/", func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusBadRequest)
		})
		r.Post("/update/counter/{name}/{value}", s.postCounter)

		r.Post("/update/*", func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusNotImplemented)
		})
	})

	return mux
}

func (s *serverRoutes) postGauge(rw http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	rawValue := chi.URLParam(r, "value")

	value, err := strconv.ParseFloat(rawValue, 64)
	if err != nil {
		log.Printf("post gauge value: %v", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	s.uc.SaveGauge(name, value)
	rw.WriteHeader(http.StatusOK)
}

func (s *serverRoutes) postCounter(rw http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	rawValue := chi.URLParam(r, "value")

	value, err := strconv.Atoi(rawValue)
	if err != nil {
		log.Printf("post counter value: %v", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	s.uc.SaveCounter(name, value)
	rw.WriteHeader(http.StatusOK)
}

func (s *serverRoutes) getListOfMetrics(rw http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./web/templates/metrics_list.html"))
	data := s.uc.GetAllMetrics()
	rw.Header().Set("Content-Type", "text/html")
	err := tmpl.Execute(rw, data)
	if err != nil {
		log.Printf("apply template: %v", err)
	}
}

func (s *serverRoutes) getMetricValue(rw http.ResponseWriter, r *http.Request) {
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

func (s *serverRoutes) postValueByJSON(rw http.ResponseWriter, r *http.Request) {
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
		log.Printf("read value from storage: %v", err)
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(respBody)
}

func (s *serverRoutes) postUpdateByJSON(rw http.ResponseWriter, r *http.Request) {
	// Обработать JSON из тела запроса - сохранить в соответствующую метрику переданное значение
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("read request body %q: %v", r.URL.Path, err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.uc.SaveValueByJSON(body)
	if err != nil {
		log.Printf("save value to storage: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

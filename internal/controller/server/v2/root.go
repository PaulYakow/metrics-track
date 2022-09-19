package v2

import (
	"fmt"
	"github.com/PaulYakow/metrics-track/internal/pkg/logger"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

const (
	pingRoute = "/ping"

	templateName = "./web/templates/metrics_list.gohtml"
)

type rootRoutes struct {
	uc     usecase.IServer
	logger logger.ILogger
}

func newRootRoutes(handler *gin.RouterGroup, uc usecase.IServer, l logger.ILogger) {
	r := &rootRoutes{
		uc:     uc,
		logger: l,
	}

	handler.GET(rootRoute, r.listOfMetrics)
	handler.GET(pingRoute, r.pingRepo)
}

func (r *rootRoutes) listOfMetrics(c *gin.Context) {
	tmpl := template.Must(template.ParseFiles(templateName))
	data, err := r.uc.GetAll(c.Request.Context())
	if err != nil {
		r.logger.Error(fmt.Errorf("router - GetAll metrics failed: %w", err))
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	c.Header("Content-Type", "text/html")
	err = tmpl.Execute(c.Writer, data)
	if err != nil {
		r.logger.Error(fmt.Errorf("router - apply template failed: %w", err))
	}
}

func (r *rootRoutes) pingRepo(c *gin.Context) {
	err := r.uc.CheckRepo()
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}

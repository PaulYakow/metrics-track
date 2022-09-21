package v2

import (
	"encoding/json"
	"fmt"
	"github.com/PaulYakow/metrics-track/internal/entity"
	"github.com/PaulYakow/metrics-track/internal/pkg/logger"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	valueRoute = "/value/"
)

type valueRoutes struct {
	uc     usecase.IServer
	logger logger.ILogger
}

func newValueRoutes(handler *gin.RouterGroup, uc usecase.IServer, l logger.ILogger) {
	v := &valueRoutes{
		uc:     uc,
		logger: l,
	}

	// get value by URL
	handler.GET(valueRoute+":type/:name",
		[]gin.HandlerFunc{
			checkContentType,
			v.bindURI,
			v.readMetricByURI,
		}...,
	)

	// get value by JSON
	handler.POST(valueRoute,
		[]gin.HandlerFunc{
			checkContentType,
			readRequestBody(v.logger),
			unmarshalJSONRequest(v.logger),
			v.readMetricByJSON,
		}...,
	)
}

func (v *valueRoutes) bindURI(c *gin.Context) {
	if c.Value(keyContentType) != valContentIsText {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var req readByURIRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.Set(keyGetURIReq, req)
	c.Next() // read metric
}

func (v *valueRoutes) readMetricByURI(c *gin.Context) {
	req, ok := c.Value(keyGetURIReq).(readByURIRequest)
	if !ok {
		v.logger.Error(fmt.Errorf("router - bad request (create): %v", req))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	metric, err := v.uc.Get(c.Request.Context(),
		entity.Metric{
			ID:    req.Name,
			MType: req.Type,
		})
	if err != nil {
		v.logger.Error(err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	respBody := []byte(metric.GetValue())
	c.Header("Content-Type", "application/json")
	c.Writer.Write(respBody)
}

func (v *valueRoutes) readMetricByJSON(c *gin.Context) {
	metric, ok := c.Value(keyUpdJSONReq).(*entity.Metric)
	if !ok {
		v.logger.Info("process single - (no) data in context key: %s", keyUpdJSONReq)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	v.logger.Info("router - POST /value: %v", *metric)

	metric, err := v.uc.Get(c.Request.Context(), *metric)
	if err != nil {
		v.logger.Error(fmt.Errorf("router - get value: %v", err))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	respBody, err := json.Marshal(&metric)
	if err != nil {
		v.logger.Error(fmt.Errorf("router - marshal to json: %v", err))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Header("Content-Type", "application/json")
	c.Writer.Write(respBody)
}

package v2

import (
	"errors"
	"fmt"
	"github.com/PaulYakow/metrics-track/internal/entity"
	"github.com/PaulYakow/metrics-track/internal/pkg/logger"
	"github.com/PaulYakow/metrics-track/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	updateRoute      = "/update/"
	batchUpdateRoute = "/updates/"
)

type updateRoutes struct {
	uc     usecase.IServer
	logger logger.ILogger
}

func newUpdateRoutes(handler *gin.RouterGroup, uc usecase.IServer, l logger.ILogger) {
	u := &updateRoutes{
		uc:     uc,
		logger: l,
	}

	// update by URL
	h := handler.Group(updateRoute)
	{
		h.POST(":type/", func(c *gin.Context) {
			c.AbortWithStatus(http.StatusNotFound)
		})

		h.POST(":type/:name/", func(c *gin.Context) {
			c.AbortWithStatus(http.StatusBadRequest)
		})

		h.POST(":type/:name/:value",
			[]gin.HandlerFunc{
				checkContentType,
				u.bindURI,
				u.createMetric,
				u.processURIRequest,
			}...,
		)
	}

	// update by single JSON
	handler.POST(updateRoute,
		[]gin.HandlerFunc{
			checkContentType,
			readRequestBody(u.logger),
			unmarshalJSONRequest(u.logger),
			u.processSingleJSONRequest,
		}...,
	)

	// update by JSON batch
	handler.POST(batchUpdateRoute,
		[]gin.HandlerFunc{
			checkContentType,
			readRequestBody(u.logger),
			unmarshalBatchRequest(u.logger),
			u.processBatchJSONRequest,
		}...,
	)
}

func (u *updateRoutes) bindURI(c *gin.Context) {
	if c.Value(keyContentType) != valContentIsText {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var req updateByURIRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.Set(keyUpdURIReq, req)
	c.Next() // create metric
}

func (u *updateRoutes) createMetric(c *gin.Context) {
	req, ok := c.Value(keyUpdURIReq).(updateByURIRequest)
	if !ok {
		u.logger.Error(fmt.Errorf("router - bad request (create): %v", req))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	metric, err := entity.Create(req.Type, req.Name, req.Value)
	if err != nil {
		u.logger.Error(fmt.Errorf("router - create temp metric: %w", err))

		if errors.Is(err, entity.ErrParseValue) {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if errors.Is(err, entity.ErrUnknownType) {
			c.AbortWithStatus(http.StatusNotImplemented)
			return
		}

		return
	}

	c.Set(keyUpdURIReq, metric)
	c.Next() // save metric
}

func (u *updateRoutes) processURIRequest(c *gin.Context) {
	metric, ok := c.Value(keyUpdURIReq).(*entity.Metric)
	if !ok {
		u.logger.Error(fmt.Errorf("router - bad request (process URL): %v", metric))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := u.uc.Save(metric); err != nil {
		u.logger.Error(fmt.Errorf("router - save metric: %w", err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
}

func (u *updateRoutes) processSingleJSONRequest(c *gin.Context) {
	metric, ok := c.Value(keyUpdJSONReq).(*entity.Metric)
	if !ok {
		u.logger.Info("process single - (no) data in context key: %s", keyUpdJSONReq)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	u.logger.Info("router - POST /update: %v", *metric)

	if err := u.uc.Save(metric); err != nil {
		u.logger.Error(fmt.Errorf("process single - save to storage: %w", err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
}

func (u *updateRoutes) processBatchJSONRequest(c *gin.Context) {
	metrics, ok := c.Value(keyUpdJSONReq).(*[]entity.Metric)
	if !ok {
		u.logger.Error(fmt.Errorf("router - bad request (process batch JSON): %v", metrics))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	u.logger.Info("router - POST /updates: %v", *metrics)

	if err := u.uc.SaveBatch(*metrics); err != nil {
		u.logger.Error(fmt.Errorf("router - batch update save to storage: %w", err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
}

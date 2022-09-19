package v2

import (
	"encoding/json"
	"fmt"
	"github.com/PaulYakow/metrics-track/internal/entity"
	"github.com/PaulYakow/metrics-track/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

const (
	keyContentType = "ContentType"    // "uri" or "json" (const values below)
	keyUpdURIReq   = "UpdUriRequest"  // update request when content is URL
	keyUpdJSONReq  = "UpdJSONRequest" // update request when content is JSON
	keyGetURIReq   = "GetUriRequest"  // get value request when content is URL
	keyGetJSONReq  = "GetJSONRequest" // get value request when content is JSON

	valContentIsText = "uri"
	valContentIsJSON = "json"
)

type updateByURIRequest struct {
	Type  string `uri:"type" binding:"required"`
	Name  string `uri:"name" binding:"required"`
	Value string `uri:"value" binding:"required"`
}

type readByURIRequest struct {
	Type string `uri:"type" binding:"required"`
	Name string `uri:"name" binding:"required"`
}

func checkContentType(c *gin.Context) {
	switch c.ContentType() {
	case "text/plain", "":
		c.Set(keyContentType, valContentIsText)
	case "application/json":
		c.Set(keyContentType, valContentIsJSON)
	default:
		c.AbortWithStatus(http.StatusNotImplemented)
		return
	}

	c.Next()
}

func readRequestBody(logger logger.ILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Value(keyContentType) != valContentIsJSON {
			logger.Error(fmt.Errorf("router - unknown content type: %q", c.ContentType()))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logger.Error(fmt.Errorf("router - read request body %q: %w", c.Request.URL.Path, err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Set(keyUpdJSONReq, body)
		c.Next()
	}
}

func unmarshalJSONRequest(logger logger.ILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var rawData entity.Metric

		if err := json.Unmarshal(c.Value(keyUpdJSONReq).([]byte), &rawData); err != nil {
			logger.Error(fmt.Errorf("router - update metric: %q", err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Set(keyUpdJSONReq, &rawData)
		c.Next()
	}
}

func unmarshalBatchRequest(logger logger.ILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var rawData []entity.Metric

		if err := json.Unmarshal(c.Value(keyUpdJSONReq).([]byte), &rawData); err != nil {
			logger.Error(fmt.Errorf("router - update metric: %q", err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Set(keyUpdJSONReq, &rawData)
		c.Next()
	}
}

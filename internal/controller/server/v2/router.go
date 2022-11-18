package v2

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	"github.com/PaulYakow/metrics-track/internal/pkg/logger"
	"github.com/PaulYakow/metrics-track/internal/usecase"
)

const (
	rootRoute = "/"
)

func NewRouter(uc usecase.IServer, l logger.ILogger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	handler := gin.New()
	handler.Use(gin.Recovery())
	//handler.Use(gin.Logger())
	handler.Use(gzip.Gzip(gzip.BestCompression))

	root := handler.Group(rootRoute)
	{
		newRootRoutes(root, uc, l)
		newUpdateRoutes(root, uc, l)
		newValueRoutes(root, uc, l)
	}

	return handler
}

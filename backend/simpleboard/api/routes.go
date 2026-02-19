package api

import (
	"github.com/gin-gonic/gin"

	"simpleboard/internal/handler"
)

func RegisterRoutes() *gin.Engine {
    r := gin.Default()

    r.GET("/api/health", handler.Health)
	r.POST("/api/register", handler.Register)

    return r
}

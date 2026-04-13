package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"simpleboard/internal/handler"
)

func RegisterRoutes() *gin.Engine {
	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
         AllowOrigins:     []string{"http://localhost:4200"},
         AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
         AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
         AllowCredentials: true,
     }))

	r.GET("/api/health", handler.Health)
	r.POST("/api/register", handler.Register)
	r.POST("/api/login", handler.Login)
	r.POST("/api/game", handler.Game)

	return r
}

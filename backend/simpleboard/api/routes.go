package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"simpleboard/internal/auth"
	"simpleboard/internal/handler"
)

func RegisterRoutes() *gin.Engine {
	r := gin.Default()

	// Auth middleware
	r.Use(auth.Middleware())

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	r.GET("/api/health", handler.Health)
	r.GET("/api/guest", handler.Guest)
	r.POST("/api/register", handler.Register)
	r.POST("/api/login", handler.Login)
	r.POST("/api/game", handler.Game)
	r.GET("/api/dashboard", handler.Dashboard)

	return r
}

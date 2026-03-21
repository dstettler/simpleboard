package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	//"simpleboard/pkg/response"
)

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

package controller

import (
	"github.com/b-open/jobbuzz/pkg/service"
	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func GetJobs(c *gin.Context) {
	jobs := service.GetJobs()
	c.JSON(200, gin.H{
		"success": true,
		"jobs":    jobs,
	})
}

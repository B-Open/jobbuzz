package controller

import (
	"net/http"

	"github.com/b-open/jobbuzz/pkg/service"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	Service *service.Service
}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (controller *Controller) GetJobs(c *gin.Context) {
	jobs, err := controller.Service.GetJobs()

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"jobs":    jobs,
	})
}

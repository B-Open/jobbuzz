package controller

import (
	"github.com/b-open/jobbuzz/pkg/service"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	Service service.Servicer
}

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (controller *Controller) GetJobs(c *gin.Context) {
	jobs, err := controller.Service.GetJobs()

	if err != nil {
		panic(err)
	}

	c.JSON(200, gin.H{
		"success": true,
		"jobs":    jobs,
	})
}

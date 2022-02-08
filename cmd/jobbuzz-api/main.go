package main

import (
	"github.com/b-open/jobbuzz/internal/config"
	"github.com/b-open/jobbuzz/pkg/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDb()
	r := gin.Default()
	r.GET("/api/jobs", controller.GetJobs)
	r.Run()
}

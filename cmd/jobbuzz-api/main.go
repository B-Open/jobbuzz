package main

import (
	"github.com/b-open/jobbuzz/internal"
	"github.com/b-open/jobbuzz/pkg/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", controller.Ping)
	logger := internal.GetLogger()
	logger.Error("Test")
	r.Run()

}

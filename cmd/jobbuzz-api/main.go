package main

import (
	"github.com/b-open/jobbuzz/pkg/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", controller.Ping)
	r.Run()
}

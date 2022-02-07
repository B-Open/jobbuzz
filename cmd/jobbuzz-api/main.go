package main

import (
	"github.com/b-open/jobbuzz/cmd/jobbuzz-api/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", controller.Ping)
	r.Run()
}

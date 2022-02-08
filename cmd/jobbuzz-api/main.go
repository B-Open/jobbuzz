package main

import (
	"github.com/b-open/jobbuzz/pkg/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", controller.Ping)
	r.GET("/scrape-jobcenter", controller.ScrapeJobcenter)
	r.GET("/scrape-bruneida", controller.ScrapeBruneida)
	r.Run()

}

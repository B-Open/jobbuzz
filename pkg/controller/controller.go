package controller

import (
	"github.com/b-open/jobbuzz/pkg/scraper"
	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func ScrapeJobcenter(c *gin.Context) {
	jobs := scraper.ScrapeJobcenter()
	c.JSON(200, gin.H{
		"success": true,
		"jobs":    jobs,
	})
}

func ScrapeBruneida(c *gin.Context) {
	jobs := scraper.ScrapeBruneida()
	c.JSON(200, gin.H{
		"success": true,
		"jobs":    jobs,
	})
}

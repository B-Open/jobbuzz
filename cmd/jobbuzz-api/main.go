package main

import (
	"log"

	"github.com/b-open/jobbuzz/internal/config"
	"github.com/b-open/jobbuzz/pkg/controller"
	"github.com/b-open/jobbuzz/pkg/service"
	"github.com/gin-gonic/gin"
)

func main() {
	configuration, err := config.LoadConfig("../../")

	if err != nil {
		log.Fatal("Fail to load db config", err)
	}

	db, err := configuration.GetDb()

	if err != nil {
		log.Fatal("Fail to get db connection", err)
	}

	service := service.Service{DB: db}
	c := controller.Controller{Service: &service}

	r := gin.Default()

	if gin.IsDebugging() {
		r.GET("/", controller.PlaygroundHandler())
	}
	r.POST("/graphql", controller.GraphqlHandler(&service))

	apiV1 := r.Group("/api/v1")
	apiV1.GET("/job", c.GetJobs)
	r.Run()
}

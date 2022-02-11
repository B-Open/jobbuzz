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
	controller := controller.Controller{Service: &service}

	r := gin.Default()

	apiV1 := r.Group("/api/v1")
	apiV1.GET("/job", controller.GetJobs)
	r.Run()
}

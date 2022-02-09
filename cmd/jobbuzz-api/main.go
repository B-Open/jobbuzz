package main

import (
	"log"

	"github.com/b-open/jobbuzz/internal/config"
	"github.com/b-open/jobbuzz/pkg/controller"
	"github.com/b-open/jobbuzz/pkg/service"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := config.GetDb()

	if err != nil {
		log.Fatal("Fail to get db connection", err)
	}

	service := service.Service{Database: db}
	controller := controller.Controller{Service: &service}

	r := gin.Default()
	r.GET("/api/jobs", controller.GetJobs)
	r.Run()
}

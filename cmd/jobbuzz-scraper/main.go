package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/b-open/jobbuzz/internal/config"
	"github.com/b-open/jobbuzz/pkg/scraper"
	"github.com/b-open/jobbuzz/pkg/service"
)

func main() {

	configuration, err := config.LoadConfig("../../")

	if err != nil {
		log.Fatal("Fail to load db config", err)
	}

	dbConfig := configuration.DbConfig

	db, err := config.GetDb(dbConfig)

	if err != nil {
		log.Fatal("Fail to get db connection", err)
	}

	service := service.Service{DB: db}

	fmt.Println("Fetching jobs from JobCenter")
	jobs := scraper.ScrapeJobcenter()

	for _, job := range jobs {
		service.CreateJob(&job)
	}

	json_jobs, err := json.Marshal(jobs)

	if err != nil {
		fmt.Println("Error json marshal", err)
	}

	fmt.Println("Printing jobs from JobCenter")
	fmt.Println(string(json_jobs))

	fmt.Println("Fetching jobs from Bruneida")
	jobs = scraper.ScrapeBruneida()

	for _, job := range jobs {
		service.CreateJob(&job)
	}

	json_jobs, err = json.Marshal(jobs)

	if err != nil {
		fmt.Println("Error json marshal", err)
	}

	fmt.Println("Printing jobs from Bruneida")
	fmt.Println(string(json_jobs))
}

package main

import (
	"encoding/json"
	"fmt"

	"github.com/b-open/jobbuzz/pkg/scraper"
)

func main() {

	// configuration, err := config.LoadConfig("../../")

	// if err != nil {
	// 	log.Fatal("Fail to load db config", err)
	// }

	// db, err := configuration.GetDb()

	// if err != nil {
	// 	log.Fatal("Fail to get db connection", err)
	// }

	// service := service.Service{DB: db}

	// Scrape JobCenter

	// fmt.Println("Fetching jobs from JobCenter")
	// jobs, _ := scraper.ScrapeJobcenter()

	// jsonJob, _ := json.Marshal(jobs)

	// fmt.Println(string(jsonJob))

	// for _, job := range jobs {
	// 	_, err = service.CreateJob(&job)

	// 	if err != nil {
	// 		fmt.Println("Fail to create job")
	// 	}
	// }

	// Scrape Bruneida

	fmt.Println("Fetching jobs from Bruneida")
	jobs, _ := scraper.ScrapeBruneida()

	jsonJob, _ := json.Marshal(jobs)

	fmt.Println(string(jsonJob))

	// for _, job := range jobs {
	// 	_, err = service.CreateJob(&job)

	// 	if err != nil {
	// 		fmt.Println("Fail to create job")
	// 	}
	// }
}

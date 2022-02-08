package main

import (
	"encoding/json"
	"fmt"

	"github.com/b-open/jobbuzz/internal/config"
	"github.com/b-open/jobbuzz/pkg/scraper"
)

func main() {
	config.InitDb()
	db := config.GetDb()

	fmt.Println("Fetching jobs from JobCenter")
	jobs := scraper.ScrapeJobcenter()

	for _, job := range jobs {
		db.Create(&job)
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
		db.Create(&job)
	}

	json_jobs, err = json.Marshal(jobs)

	if err != nil {
		fmt.Println("Error json marshal", err)
	}

	fmt.Println("Printing jobs from Bruneida")
	fmt.Println(string(json_jobs))
}

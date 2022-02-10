package main

import (
	"encoding/json"
	"fmt"

	"github.com/b-open/jobbuzz/pkg/scraper"
)

func main() {
	// This is just for me to test the scripts
	jobs := scraper.ScrapeBruneida()

	json_jobs, err := json.Marshal(jobs)

	if err != nil {
		fmt.Println("Error json marshal", err)
	}

	fmt.Println("Printing jobs from JobCenter")
	fmt.Println(string(json_jobs))
}

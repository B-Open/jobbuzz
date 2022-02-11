package scraper

import (
	"fmt"

	"github.com/b-open/jobbuzz/pkg/model"
	"github.com/gocolly/colly"
)

const (
	JobCenter = 1
	Bruneida  = 2
)

func HandleError(collectors []*colly.Collector) {
	for _, collector := range collectors {
		collector.OnError(func(r *colly.Response, err error) {
			fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
		})
	}
}

func HandleRequest(collectors []*colly.Collector) {
	for _, collector := range collectors {
		collector.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL.String())
		})
	}
}

func ConvertJobMapToJobSlice(jobMap map[string]model.Job) []model.Job {
	jobs := []model.Job{}
	for _, job := range jobMap {
		jobs = append(jobs, job)
	}

	return jobs
}

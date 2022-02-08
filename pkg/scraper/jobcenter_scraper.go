package scraper

import (
	"fmt"
	"strconv"

	"github.com/b-open/jobbuzz/pkg/model"
	"github.com/gocolly/colly"
)

func ScrapeJobcenter() []model.Job {
	jobs := []model.Job{}

	c := colly.NewCollector(
		colly.AllowedDomains("www.jobcentrebrunei.gov.bn"),
	)

	c.OnHTML("li.list-group-item.list-group-item-flex", func(e *colly.HTMLElement) {
		job_title := e.ChildText(".jp_job_post_right_cont h4 a")
		company := e.ChildText(".jp_job_post_right_cont p a")
		salary := e.ChildText(".jp_job_post_right_cont>ul li:first-child")
		location := e.ChildText(".jp_job_post_right_cont>ul li:nth-child(2)")

		job := model.Job{
			Title:    job_title,
			Company:  company,
			Salary:   salary,
			Location: location,
		}

		jobs = append(jobs, job)

	})

	collectors := []*colly.Collector{c}

	HandleError(collectors)

	HandleRequest(collectors)

	// Limit to two pages
	for i := 1; i < 3; i++ {

		url := fmt.Sprintf("https://www.jobcentrebrunei.gov.bn/web/guest/search-job?q=&delta=200&start=%s", strconv.Itoa(i))

		if err := c.Visit(url); err != nil {
			fmt.Println("Error: ", err)
			break
		}
	}

	return jobs
}
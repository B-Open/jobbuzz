package scraper

import (
	"fmt"
	"strconv"

	"github.com/b-open/jobbuzz/pkg/types"
	"github.com/gocolly/colly"
)

func ScrapeJobcenter() types.Jobs {
	jobs := types.Jobs{}

	c := colly.NewCollector(
		colly.AllowedDomains("www.jobcentrebrunei.gov.bn"),
	)

	c.OnHTML("li.list-group-item.list-group-item-flex", func(e *colly.HTMLElement) {
		job_title := e.ChildText(".jp_job_post_right_cont h4 a")
		company := e.ChildText(".jp_job_post_right_cont p a")
		salary := e.ChildText(".jp_job_post_right_cont>ul li:first-child")
		location := e.ChildText(".jp_job_post_right_cont>ul li:nth-child(2)")

		job := types.Job{
			Title:    job_title,
			Company:  company,
			Salary:   salary,
			Location: location,
		}

		jobs = append(jobs, job)

	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Limit to two pages
	for i := 1; i < 3; i++ {

		url := fmt.Sprintf("https://www.jobcentrebrunei.gov.bn/web/guest/search-job?q=&start=%s", strconv.Itoa(i))

		if err := c.Visit(url); err != nil {
			fmt.Println("Error: ", err)
			break
		}
	}

	return jobs
}

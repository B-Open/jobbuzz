package scraper

import (
	"fmt"
	"strconv"

	"github.com/b-open/jobbuzz/internal/util"
	"github.com/b-open/jobbuzz/pkg/model"
	"github.com/gocolly/colly"
)

const (
	delta = 10
)

func ScrapeJobcenter() []model.Job {

	jobs := []model.Job{}
	descriptions := []string{}

	linksCollector := colly.NewCollector(
		colly.AllowedDomains("www.jobcentrebrunei.gov.bn"),
	)

	jobsCollector := linksCollector.Clone()

	linksCollector.OnHTML("li.list-group-item.list-group-item-flex", func(e *colly.HTMLElement) {
		job_title := e.ChildText(".jp_job_post_right_cont h4 a")
		company := e.ChildText(".jp_job_post_right_cont p a")
		salary := e.ChildText(".jp_job_post_right_cont>ul li:first-child")
		location := e.ChildText(".jp_job_post_right_cont>ul li:nth-child(2)")
		link := e.ChildAttr(".jp_job_post_right_cont h4 a", "href")

		job := model.Job{
			Title:    job_title,
			Company:  company,
			Salary:   salary,
			Location: location,
			Link:     link,
		}

		jobUrl := fmt.Sprintf("https://www.jobcentrebrunei.gov.bn%s", link)

		if err := jobsCollector.Visit(jobUrl); err != nil {
			fmt.Println(jobUrl)
			fmt.Println("Error: ", err)
		}

		jobs = append(jobs, job)

	})

	jobsCollector.OnHTML("body", func(h *colly.HTMLElement) {
		description := h.ChildText(".container .row .col-lg-8.col-md-12.col-sm-12.col-12")

		description = util.StandardizeSpaces(description)

		descriptions = append(descriptions, description)
	})

	collectors := []*colly.Collector{linksCollector, jobsCollector}

	HandleError(collectors)

	HandleRequest(collectors)

	// Limit to two pages
	for i := 1; i < 3; i++ {

		url := fmt.Sprintf("https://www.jobcentrebrunei.gov.bn/web/guest/search-job?q=&delta=%s&start=%s", strconv.Itoa(delta), strconv.Itoa(i))

		if err := linksCollector.Visit(url); err != nil {
			fmt.Println("Error: ", err)
			break
		}
	}

	for i := range jobs {
		jobs[i].Description = descriptions[i]
	}

	return jobs
}

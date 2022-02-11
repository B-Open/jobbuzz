package scraper

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/b-open/jobbuzz/internal/util"
	"github.com/b-open/jobbuzz/pkg/model"
	"github.com/gocolly/colly"
)

const (
	pageSize          = 30
	jobcenterProvider = "jobcenter"
)

func ScrapeJobcenter() []model.Job {

	jobMap := map[string]model.Job{}

	linksCollector := colly.NewCollector(
		colly.AllowedDomains("www.jobcentrebrunei.gov.bn"),
	)

	jobsCollector := linksCollector.Clone()

	linksCollector.OnHTML("body", func(e *colly.HTMLElement) {
		e.ForEachWithBreak("li.list-group-item.list-group-item-flex", func(i int, h *colly.HTMLElement) bool {
			jobTitle := h.ChildText(".jp_job_post_right_cont h4 a")
			company := h.ChildText(".jp_job_post_right_cont p a")
			salary := h.ChildText(".jp_job_post_right_cont>ul li:first-child")
			location := h.ChildText(".jp_job_post_right_cont>ul li:nth-child(2)")
			link := h.ChildAttr(".jp_job_post_right_cont h4 a", "href")

			jobId, err := getJobcenterJobId(link)

			// This will continue the loop
			if err != nil {
				return true
			}

			job := model.Job{
				Provider: jobcenterProvider,
				JobId:    jobId,
				Title:    jobTitle,
				Company:  company,
				Salary:   salary,
				Location: location,
				Link:     link,
			}

			jobUrl := fmt.Sprintf("https://www.jobcentrebrunei.gov.bn%s", link)

			if err := jobsCollector.Visit(jobUrl); err != nil {
				fmt.Println(jobUrl)
				fmt.Println("Error: ", err)

				return true
			}

			jobMap[jobId] = job

			return true
		})

	})

	jobsCollector.OnHTML("body", func(h *colly.HTMLElement) {
		link := h.Request.URL.String()
		jobId, err := getJobcenterJobId(link)

		if err == nil {
			description := h.ChildText(".container .row .col-lg-8.col-md-12.col-sm-12.col-12")

			description = util.StandardizeSpaces(description)

			job := jobMap[jobId]
			job.Description = description
			jobMap[jobId] = job
		}

	})

	collectors := []*colly.Collector{linksCollector, jobsCollector}

	HandleError(collectors)

	HandleRequest(collectors)

	// Limit to two pages
	for i := 1; i < 3; i++ {

		url := fmt.Sprintf("https://www.jobcentrebrunei.gov.bn/web/guest/search-job?q=&delta=%d&start=%d", pageSize, i)

		if err := linksCollector.Visit(url); err != nil {
			fmt.Println("Error: ", err)
			break
		}
	}

	jobs := ConvertJobMapToJobSlice(jobMap)

	return jobs
}

func getJobcenterJobId(s string) (string, error) {
	r := regexp.MustCompile(`^\/web\/guest\/view-job\/-\/jobs\/(?P<jobId>\d+)\/.*$`)
	matches := r.FindStringSubmatch(s)

	if len(matches) < 1 {
		return "", errors.New("no job id found")
	}

	if matches[1] == "" {
		return "", errors.New("no job id found")
	}

	return matches[1], nil
}

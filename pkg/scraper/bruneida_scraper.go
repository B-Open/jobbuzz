package scraper

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/b-open/jobbuzz/internal/util"
	"github.com/b-open/jobbuzz/pkg/model"
	"github.com/gocolly/colly"
)

const (
	bruneidaProvider = "bruneida"
)

func ScrapeBruneida() []model.Job {
	allowedDomain := colly.AllowedDomains("www.bruneida.com")
	linkCollector := colly.NewCollector(
		allowedDomain,
		colly.Async(true),
	)

	linkCollector.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 10})

	jobCollector := linkCollector.Clone()

	jobMap := map[string]model.Job{}

	// Scraping the links
	linkCollector.OnHTML("ul.list-az.ul-azs", func(e *colly.HTMLElement) {
		e.ForEachWithBreak(".az-detail>h3.az-title>a.h-elips", func(i int, child *colly.HTMLElement) bool {
			link := child.Attr("href")

			jobId, err := getBruneidaJobId(link)

			// This will continue the loop
			if err != nil {
				return true
			}

			job := model.Job{Provider: bruneidaProvider, JobId: jobId}
			jobMap[jobId] = job

			jobCollector.Visit(link)

			return true
		})
		jobCollector.Wait()
	})

	// Scraping the jobs
	jobCollector.OnHTML("body", func(h *colly.HTMLElement) {

		jobTitle := h.ChildText("#title-box-inner div.inline-block.pull-left h1")
		company := h.ChildText("#ad-contact ul li:first-child span.bb b.small")
		salary := h.ChildText("#ad-body-inner .opt .opt-dl:nth-child(3) .dd")

		// TODO: Use HTML Minifier and whitelist
		description := util.StandardizeSpaces(h.ChildText("#full-description"))

		locations := []string{}
		h.ForEach("#ad-body-inner .opt .opt-dl", func(i int, child *colly.HTMLElement) {
			title := child.ChildText(".dt")

			if strings.Contains(title, "City") || strings.Contains(title, "Local") {
				locations = append(locations, child.ChildText(".dd"))
			}
		})

		location := strings.Join(locations, " ")

		link := h.Request.URL.String()

		jobId, err := getBruneidaJobId(link)

		if err == nil {
			job := jobMap[jobId]

			newJob := model.Job{
				JobId:       job.JobId,
				Provider:    job.Provider,
				Title:       jobTitle,
				Company:     company,
				Salary:      salary,
				Location:    location,
				Description: description,
			}

			jobMap[jobId] = newJob
		}

	})

	collectors := []*colly.Collector{linkCollector, jobCollector}

	HandleError(collectors)

	HandleRequest(collectors)

	// Limit to one page
	for i := 1; i < 2; i++ {
		url := fmt.Sprintf("https://www.bruneida.com/brunei/jobs/?&page=%s", strconv.Itoa(i))
		linkCollector.Visit(url)
	}
	linkCollector.Wait()

	jobs := ConvertJobMapToJobSlice(jobMap)

	return jobs

}

func getBruneidaJobId(s string) (string, error) {
	r := regexp.MustCompile(`-(?P<jobId>\d+)`)
	matches := r.FindStringSubmatch(s)

	if len(matches) < 1 {
		return "", errors.New("no job id found")
	}

	if matches[1] == "" {
		return "", errors.New("no job id found")
	}

	return matches[1], nil
}

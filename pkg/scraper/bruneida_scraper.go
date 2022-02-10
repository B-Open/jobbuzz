package scraper

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/b-open/jobbuzz/internal/util"
	"github.com/b-open/jobbuzz/pkg/model"
	"github.com/gocolly/colly"
)

func ScrapeBruneida() []model.Job {
	allowedDomain := colly.AllowedDomains("www.bruneida.com")
	linkCollector := colly.NewCollector(
		allowedDomain,
	)
	jobCollector := linkCollector.Clone()

	jobs := []model.Job{}
	links := []string{}
	jobIds := []string{}

	// Scraping the links
	linkCollector.OnHTML("ul.list-az.ul-azs", func(e *colly.HTMLElement) {
		e.ForEach(".az-detail>h3.az-title>a.h-elips", func(i int, child *colly.HTMLElement) {
			link := child.Attr("href")

			jobCollector.Visit(link)
			links = append(links, link)

			jobId := getBruneidaJobId(link)
			jobIds = append(jobIds, jobId)
		})
	})

	// Scraping the jobs
	jobCollector.OnHTML("body", func(h *colly.HTMLElement) {
		jobTitle := h.ChildText("#title-box-inner div.inline-block.pull-left h1")
		company := h.ChildText("#ad-contact ul li:first-child span.bb b.small")
		salary := h.ChildText("#ad-body-inner .opt .opt-dl:nth-child(3) .dd")
		description := util.StandardizeSpaces(h.ChildText("#full-description"))

		locations := []string{}
		h.ForEach("#ad-body-inner .opt .opt-dl", func(i int, child *colly.HTMLElement) {
			title := child.ChildText(".dt")

			if strings.Contains(title, "City") || strings.Contains(title, "Local") {
				locations = append(locations, child.ChildText(".dd"))
			}
		})

		location := strings.Join(locations, " ")

		job := model.Job{
			Title:       jobTitle,
			Company:     company,
			Salary:      salary,
			Location:    location,
			Description: description,
		}

		jobs = append(jobs, job)
	})

	collectors := []*colly.Collector{linkCollector, jobCollector}

	HandleError(collectors)

	HandleRequest(collectors)

	// Limit to one page
	for i := 1; i < 2; i++ {
		url := fmt.Sprintf("https://www.bruneida.com/brunei/jobs/?&page=%s", strconv.Itoa(i))
		linkCollector.Visit(url)
	}

	// Adding links to the existing jobs slice
	for i := range jobs {
		jobs[i].Link = links[i]
		jobs[i].JobId = jobIds[i]
	}

	return jobs

}

func getBruneidaJobId(s string) string {
	r := regexp.MustCompile(`-(?P<jobId>\d+)`)
	return fmt.Sprintf("bruneida-%s", r.FindStringSubmatch(s)[1])
}

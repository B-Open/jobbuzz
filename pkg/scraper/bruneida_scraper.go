package scraper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/b-open/jobbuzz/pkg/types"
	"github.com/gocolly/colly"
)

func ScrapeBruneida() types.Jobs {
	allowedDomain := colly.AllowedDomains("www.bruneida.com")
	linkCollector := colly.NewCollector(
		allowedDomain,
	)
	jobCollector := colly.NewCollector(
		allowedDomain,
	)

	// links := make(links, 0)
	jobs := make(types.Jobs, 0)

	// Scraping the links
	linkCollector.OnHTML("ul.list-az.ul-azs", func(e *colly.HTMLElement) {
		e.ForEach(".az-detail>h3.az-title>a.h-elips", func(i int, child *colly.HTMLElement) {
			link := child.Attr("href")
			jobCollector.Visit(link)

		})
	})

	// Scraping the jobs
	jobCollector.OnHTML("body", func(h *colly.HTMLElement) {
		job_title := h.ChildText("#title-box-inner div.inline-block.pull-left h1")
		company := h.ChildText("#ad-contact ul li:first-child span.bb b.small")
		salary := h.ChildText("#ad-body-inner .opt .opt-dl:nth-child(3) .dd")

		fmt.Println(job_title)
		fmt.Println(company)
		fmt.Println(salary)

		location := ""
		h.ForEach("#ad-body-inner .opt .opt-dl", func(i int, child *colly.HTMLElement) {
			title := child.ChildText(".dt")

			if strings.Contains(title, "City") || strings.Contains(title, "Local") {
				s := child.ChildText(".dd")
				location = location + s + " "
			}
		})

		job := types.Job{
			Title:    job_title,
			Company:  company,
			Salary:   salary,
			Location: location,
		}

		jobs = append(jobs, job)
	})

	// On Request Hooks
	linkCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	jobCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Limit to one page
	for i := 1; i < 2; i++ {
		url := fmt.Sprintf("https://www.bruneida.com/brunei/jobs/?&page=%s", strconv.Itoa(i))
		linkCollector.Visit(url)
	}

	jobCollector.Visit("https://www.bruneida.com/FEMALE-SALES-ASSISTANT-106578")

	return jobs

}

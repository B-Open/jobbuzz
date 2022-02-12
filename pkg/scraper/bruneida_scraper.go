package scraper

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/b-open/jobbuzz/pkg/model"
)

func ScrapeBruneida() ([]*model.Job, error) {
	jobs := []*model.Job{}

	for i := 1; i < 2; i++ {
		url := fmt.Sprintf("https://www.bruneida.com/brunei/jobs/?&page=%d", i)

		links, err := getJobLinks(url)

		if err != nil {
			return nil, err
		}

		for _, link := range links {
			job, err := scrapeBruneidaJob(link)

			if err != nil {
				fmt.Printf("Fail to scrape job for link : %s, err: %s \n", link, err)
				continue
			}

			jobs = append(jobs, job)
		}
	}

	return jobs, nil
}

func scrapeBruneidaJob(url string) (*model.Job, error) {

	doc, err := getDocument(url)

	if err != nil {
		return nil, err
	}

	jobTitle := doc.Find("#title-box-inner div.inline-block.pull-left h1").Text()
	company := doc.Find("#ad-contact ul li:first-child span.bb b.small").Text()
	salary := doc.Find("#ad-body-inner .opt .opt-dl:nth-child(3) .dd").Text()

	description, err := minifyHtml(doc.Find("#full-description").Text())

	if err != nil {
		return nil, err
	}

	locations := []string{}
	doc.Find("#ad-body-inner .opt .opt-dl").EachWithBreak(func(i int, s *goquery.Selection) bool {
		title := s.Find(".dt").Text()

		if strings.Contains(title, "City") || strings.Contains(title, "Local") {
			locations = append(locations, s.Find(".dd").Text())
		}

		return true
	})

	location := strings.Join(locations, " ")

	jobId, err := getBruneidaJobId(url)

	if err != nil {
		return nil, err
	}

	job := model.Job{
		JobId:       jobId,
		Provider:    Bruneida,
		Title:       jobTitle,
		Company:     company,
		Salary:      salary,
		Location:    location,
		Description: *description,
	}

	return &job, nil

}

func getJobLinks(url string) ([]string, error) {
	links := []string{}
	doc, err := getDocument(url)

	if err != nil {
		return nil, err
	}

	doc.Find(".az-detail>h3.az-title>a.h-elips").EachWithBreak(func(i int, s *goquery.Selection) bool {
		link, exist := s.Attr("href")

		if !exist {
			return true
		}

		links = append(links, link)

		return true
	})

	return links, nil

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

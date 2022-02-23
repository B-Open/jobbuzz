package scraper

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/b-open/jobbuzz/pkg/model"
)

const (
	pageSize = 30

	jobcenterUrl = "https://www.jobcentrebrunei.gov.bn"
)

func ScrapeJobcenter() ([]*model.Job, error) {

	jobs := []*model.Job{}

	url := fmt.Sprintf("%s/web/guest/search-job?q=&delta=%d", jobcenterUrl, pageSize)

	doc, err := getDocument(url)
	if err != nil {
		return nil, err
	}

	doc.Find("li.list-group-item.list-group-item-flex").EachWithBreak(func(i int, s *goquery.Selection) bool {

		jobTitle := s.Find(".jp_job_post_right_cont h4 a").Text()
		company := s.Find(".jp_job_post_right_cont p a").Text()
		salary := s.Find(".jp_job_post_right_cont>ul li:first-child").Text()
		location := s.Find(".jp_job_post_right_cont>ul li:nth-child(2)").Text()

		link, exist := s.Find(".jp_job_post_right_cont h4 a").Attr("href")
		if !exist {
			return true
		}

		providerJobId, err := getJobcenterJobId(link)
		if err != nil {
			return true
		}

		description, err := scrapeJobDescription(link)
		if err != nil {
			return true
		}

		job := model.Job{
			Provider:      JobCenter,
			ProviderJobID: providerJobId,
			Title:         jobTitle,
			Company:       company,
			Salary:        salary,
			Location:      location,
			Link:          link,
			Description:   *description,
		}

		jobs = append(jobs, &job)

		return true
	})

	return jobs, nil
}

func scrapeJobDescription(jobUrl string) (*string, error) {

	url := fmt.Sprintf("%s%s", jobcenterUrl, jobUrl)

	doc, err := getDocument(url)
	if err != nil {
		return nil, err
	}

	description := doc.Find(".container .row .col-lg-8.col-md-12.col-sm-12.col-12").Text()

	minifiedDescription, err := minifyHtml(description)

	if err != nil {
		return nil, err
	}

	return minifiedDescription, nil
}

func getJobcenterJobId(s string) (string, error) {

	r := regexp.MustCompile(`^\/web\/guest\/view-job\/-\/jobs\/(?P<jobId>\d+)\/.*$`)

	matches := r.FindStringSubmatch(s)

	if len(matches) < 2 {
		return "", errors.New(fmt.Sprintf("job id is empty: %s", s))
	}

	if matches[1] == "" {
		return "", errors.New(fmt.Sprintf("no job id found: %s", s))
	}

	return matches[1], nil
}

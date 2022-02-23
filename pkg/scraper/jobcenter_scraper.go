package scraper

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/b-open/jobbuzz/pkg/model"
	"github.com/rs/zerolog/log"
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
		companyName := s.Find(".jp_job_post_right_cont p a").Text()
		salary := s.Find(".jp_job_post_right_cont>ul li:first-child").Text()
		location := s.Find(".jp_job_post_right_cont>ul li:nth-child(2)").Text()

		link, exist := s.Find(".jp_job_post_right_cont h4 a").Attr("href")
		if !exist {
			return true
		}
		companyLink, companyExist := s.Find(".jp_job_post_right_cont p a").Attr("href")

		providerJobId, err := getJobcenterJobId(link)
		if err != nil {
			log.Error().Err(err)
			return true
		}
		var providerCompanyId string
		if companyExist {
			providerCompanyId, err = getJobcenterCompanyId(companyLink)
			if err != nil {
				log.Error().Err(err)
				return true
			}
		}

		description, err := scrapeJobDescription(link)
		if err != nil {
			log.Error().Err(err)
			return true
		}

		// TODO: this get created automatically but it is not efficient to have duplicate companies in all the jobs. Think of a way to only insert companies once, then attach the company id to the jobs, then insert the jobs
		var company *model.Company
		if companyExist {
			company = &model.Company{
				Provider:          JobCenter,
				ProviderCompanyID: providerCompanyId,
				Name:              companyName,
			}
		}

		job := model.Job{
			Provider:      JobCenter,
			ProviderJobID: providerJobId,
			Title:         jobTitle,
			Salary:        salary,
			Location:      location,
			Link:          link,
			Description:   *description,
			Company:       company,
		}
		jobs = append(jobs, &job)

		return true
	})

	// TODO: scrape company page

	// TODO: attach company information

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
	return getJobcenterId(s, "job", "jobs")
}

func getJobcenterCompanyId(s string) (string, error) {
	return getJobcenterId(s, "employer", "employer")
}

func getJobcenterId(url, idType1, idType2 string) (string, error) {
	pattern := fmt.Sprintf(`^\/web\/guest\/view-%s\/-\/%s\/(?P<jobId>\d+).*`, idType1, idType2)
	r := regexp.MustCompile(pattern)

	matches := r.FindStringSubmatch(url)

	if len(matches) < 2 {
		return "", errors.New(fmt.Sprintf("job id is empty: %s", url))
	}

	if matches[1] == "" {
		return "", errors.New(fmt.Sprintf("no job id found: %s", url))
	}

	return matches[1], nil
}

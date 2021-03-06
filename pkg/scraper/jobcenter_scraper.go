package scraper

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/b-open/jobbuzz/pkg/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	pageSize = 200
)

func NewJobCentreScraper() JobCentreScraper {
	return JobCentreScraper{
		BaseURL:     "https://www.jobcentrebrunei.gov.bn",
		FetchClient: &FetchClient{},
	}
}

func (s *JobCentreScraper) ScrapeJobs() ([]*model.Job, map[string]*model.Company, error) {
	lastPageNo, err := s.scrapeJobcenterLastPageNumber()
	if err != nil {
		return nil, nil, err
	}

	// get job and company links and IDs
	jobs, companies := s.scrapeJobAndCompanies(lastPageNo)

	// get job details
	jobs = s.scrapeJobDetails(jobs)

	// get company details and update company id in jobs
	companies = s.scrapeCompanyDetails(companies)

	return jobs, companies, nil
}

func (s *JobCentreScraper) scrapeJobAndCompanies(maxPage int) ([]*model.Job, map[string]*model.Company) {
	var wg sync.WaitGroup
	var jobs []*model.Job
	var companies = map[string]*model.Company{}

	for i := 1; i <= maxPage; i++ {
		wg.Add(1)

		go func(page int) {
			defer wg.Done()
			url := fmt.Sprintf("%s/web/guest/search-job?q=&delta=%d&start=%d", s.BaseURL, pageSize, page)

			doc, err := s.FetchClient.GetDocument(url)
			if err != nil {
				log.Error().Err(err).Msgf("Fail to scrape url : %s", url)
				return
			}

			// get job and company links and IDs
			doc.Find("li.list-group-item.list-group-item-flex").Each(func(i int, s *goquery.Selection) {
				jobTitle := s.Find(".jp_job_post_right_cont h4 a").Text()
				companyName := s.Find(".jp_job_post_right_cont p a").Text()
				salary := s.Find(".jp_job_post_right_cont>ul li:first-child").Text()
				location := s.Find(".jp_job_post_right_cont>ul li:nth-child(2)").Text()

				link, exist := s.Find(".jp_job_post_right_cont h4 a").Attr("href")
				if !exist {
					return
				}
				companyLink, companyExist := s.Find(".jp_job_post_right_cont p a").Attr("href")
				providerJobId, err := getJobcenterJobId(link)
				if err != nil {
					log.Warn().Err(err)
					return
				}
				var providerCompanyId string
				if companyExist {
					providerCompanyId, err = getJobcenterCompanyId(companyLink)
					if err != nil {
						log.Warn().Err(err)
						return
					}
				}

				var company *model.Company
				if companyExist {
					company = &model.Company{
						Provider:          JobCenter,
						ProviderCompanyID: providerCompanyId,
						Name:              companyName,
						Link:              companyLink,
					}
					if _, ok := companies[providerCompanyId]; !ok {
						companies[providerCompanyId] = company
					}
				}

				job := model.Job{
					Provider:      JobCenter,
					ProviderJobID: providerJobId,
					Link:          link,
					Company:       company,
					Title:         jobTitle,
					Salary:        salary,
					Location:      location,
				}
				jobs = append(jobs, &job)
			})
		}(i)
	}

	wg.Wait()

	return jobs, companies
}

func (s *JobCentreScraper) scrapeJobDetails(jobs []*model.Job) []*model.Job {
	var wg sync.WaitGroup

	for _, job := range jobs {
		wg.Add(1)

		go func(job *model.Job) {
			defer wg.Done()
			doc, err := s.FetchClient.GetDocument(s.BaseURL + job.Link)
			if err != nil {
				log.Error().Err(err).Msgf("Fail to scrape url : %s", job.Link)
				return
			}

			doc.Find("li.list-group-item.list-group-item-flex").Each(func(i int, s *goquery.Selection) {
				description := doc.Find(".container .row .col-lg-8.col-md-12.col-sm-12.col-12").Text()

				minifiedDescription, err := minifyHtml(description)
				if err != nil {
					log.Warn().Err(err)
					return
				}

				job.Description = *minifiedDescription
			})
		}(job)
	}

	wg.Wait()

	return jobs
}

func (s *JobCentreScraper) scrapeCompanyDetails(companies map[string]*model.Company) map[string]*model.Company {
	var wg sync.WaitGroup

	for _, company := range companies {
		wg.Add(1)

		go func(company *model.Company) {
			defer wg.Done()
			doc, err := s.FetchClient.GetDocument(s.BaseURL + company.Link)
			if err != nil {
				log.Error().Err(err).Msgf("Fail to scrape url : %s", company.Link)
				return
			}

			registrationNo := doc.Find("#collapseTwentyLeftFour > div:nth-child(1) > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(2) > td:nth-child(3)").Text()
			description := doc.Find(".jp_job_des > p:nth-child(2)").Text()
			minifiedDescription, err := minifyHtml(description)
			if err != nil {
				log.Warn().Err(err)
				return
			}

			// TODO: change to a proper struct
			// TODO: add more data
			content := map[string]interface{}{
				"RegistrationNo": registrationNo,
				"Description":    minifiedDescription,
			}

			contentBytes, err := json.Marshal(content)
			if err != nil {
				log.Warn().Err(err).Msgf("Error marshalling content to json: %+v", content)
			}

			company.Content = string(contentBytes)
		}(company)
	}

	wg.Wait()

	return companies
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
		return "", fmt.Errorf("job id is empty: %s", url)
	}

	if matches[1] == "" {
		return "", fmt.Errorf("no job id found: %s", url)
	}

	return matches[1], nil
}

func (s *JobCentreScraper) scrapeJobcenterLastPageNumber() (int, error) {

	urlString := fmt.Sprintf("%s/web/guest/search-job?q=&delta=%d&start=%d", s.BaseURL, pageSize, 1)

	doc, err := s.FetchClient.GetDocument(urlString)

	if err != nil {
		return 0, errors.Wrapf(err, "failed to get document. url: %s", urlString)
	}

	pageNoAsString := doc.Find("ul.pagination>li.page-item:nth-last-child(2)>a").Text()

	pageNoAsString = strings.ReplaceAll(pageNoAsString, "Page", "")

	pageNo, err := strconv.Atoi(pageNoAsString)

	if err != nil {
		return 0, errors.Wrapf(err, "failed to get page number. url: %s", urlString)
	}

	return pageNo, nil
}

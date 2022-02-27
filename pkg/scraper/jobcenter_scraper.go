package scraper

import (
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
	pageSize     = 200
	jobcenterUrl = "https://www.jobcentrebrunei.gov.bn"
)

func NewJobCentreScraper() JobCentreScraper {
	return JobCentreScraper{FetchClient: &FetchClient{}}
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
			url := fmt.Sprintf("%s/web/guest/search-job?q=&delta=%d&start=%d", jobcenterUrl, pageSize, page)

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
			doc, err := s.FetchClient.GetDocument(job.Link)
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

	urlString := fmt.Sprintf("%s/web/guest/search-job?q=&delta=%d&start=%d", jobcenterUrl, pageSize, 1)

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

package scraper

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/b-open/jobbuzz/pkg/model"
	"github.com/rs/zerolog/log"
	"github.com/pkg/errors"
)

const (
	pageSize = 200

	jobcenterUrl = "https://www.jobcentrebrunei.gov.bn"
)

func ScrapeJobcenter() ([]*model.Job, error) {

	jobcenterScraper := createScraper()

	lastPageNo, err := scrapeJobcenterLastPageNumber()
	if err != nil {
		return nil, err
	}

	for i := 1; i <= lastPageNo; i++ {
		jobcenterScraper.wg.Add(1)
		url := fmt.Sprintf("%s/web/guest/search-job?q=&delta=%d&start=%d", jobcenterUrl, pageSize, i)
		go (&jobcenterScraper).scrapeJobcenterJobsListing(url)
	}

	jobcenterScraper.wg.Wait()

	return jobcenterScraper.jobs, nil
}

func (jobcenterScraper *scraper) scrapeJobcenterJobsListing(url string) bool {
	defer jobcenterScraper.wg.Done()

	doc, err := getDocument(url)
	if err != nil {
		fmt.Printf("Fail to scrape url : %s, err: %s \n", url, err)
		return true
	}

	doc.Find("li.list-group-item.list-group-item-flex").Each(func(i int, s *goquery.Selection) {
		jobcenterScraper.wg.Add(1)

		go jobcenterScraper.scrapeJobcenterJob(s)
	})

	return true
}

func (jobcenterScraper *scraper) scrapeJobcenterJob(s *goquery.Selection) bool {
	defer jobcenterScraper.wg.Done()
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
	jobcenterScraper.jobs = append(jobcenterScraper.jobs, &job)

	return true

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
		return "", fmt.Errorf("job id is empty: %s", url)
	}

	if matches[1] == "" {
		return "", fmt.Errorf("no job id found: %s", url)
	}

	return matches[1], nil
}

func scrapeJobcenterLastPageNumber() (int, error) {

	urlString := fmt.Sprintf("%s/web/guest/search-job?q=&delta=%d&start=%d", jobcenterUrl, pageSize, 1)

	doc, err := getDocument(urlString)

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

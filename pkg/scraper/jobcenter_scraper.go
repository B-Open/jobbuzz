package scraper

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/b-open/jobbuzz/pkg/model"
)

const (
	pageSize = 200

	jobcenterUrl = "https://www.jobcentrebrunei.gov.bn"
)

var jobcenterJobs = []*model.Job{}

func ScrapeJobcenter() ([]*model.Job, error) {
	var wg sync.WaitGroup

	lastPageNo, err := scrapeJobcenterLastPageNumber()
	if err != nil {
		return nil, err
	}

	for i := 1; i <= lastPageNo; i++ {
		wg.Add(1)
		url := fmt.Sprintf("%s/web/guest/search-job?q=&delta=%d&start=%d", jobcenterUrl, pageSize, i)
		go scrapeJobcenterJobsListing(url, &wg)
	}

	wg.Wait()

	return jobcenterJobs, nil
}

func scrapeJobcenterJobsListing(url string, wg *sync.WaitGroup) {
	defer wg.Done()

	doc, err := getDocument(url)
	if err != nil {
		fmt.Printf("Fail to scrape url : %s, err: %s \n", url, err)
	}

	doc.Find("li.list-group-item.list-group-item-flex").Each(func(i int, s *goquery.Selection) {
		wg.Add(1)

		go func(s *goquery.Selection) {
			defer wg.Done()
			jobTitle := s.Find(".jp_job_post_right_cont h4 a").Text()
			company := s.Find(".jp_job_post_right_cont p a").Text()
			salary := s.Find(".jp_job_post_right_cont>ul li:first-child").Text()
			location := s.Find(".jp_job_post_right_cont>ul li:nth-child(2)").Text()

			link, exist := s.Find(".jp_job_post_right_cont h4 a").Attr("href")
			if !exist {
				fmt.Printf("Fail to scrape job link : %s, err: %s \n", link, err)
			}

			providerJobId, err := getJobcenterJobId(link)
			if err != nil {
				fmt.Printf("Fail to scrape job id for link : %s, err: %s \n", link, err)
			}

			description, err := scrapeJobDescription(link)
			if err != nil {
				fmt.Printf("Fail to scrape job description : %s, err: %s \n", link, err)
			}

			job := model.Job{
				Provider:      JobCenter,
				ProviderJobId: providerJobId,
				Title:         jobTitle,
				Company:       company,
				Salary:        salary,
				Location:      location,
				Link:          link,
				Description:   *description,
			}

			jobcenterJobs = append(jobcenterJobs, &job)

		}(s)
	})
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
		return "", fmt.Errorf("job id is empty: %s", s)
	}

	if matches[1] == "" {
		return "", fmt.Errorf("no job id found: %s", s)
	}

	return matches[1], nil
}

func scrapeJobcenterLastPageNumber() (int, error) {

	urlString := fmt.Sprintf("%s/web/guest/search-job?q=&delta=%d&start=%d", jobcenterUrl, pageSize, 1)

	doc, err := getDocument(urlString)

	if err != nil {
		return 0, errors.New("failed to get document")
	}

	pageNoAsString := doc.Find("ul.pagination>li.page-item:nth-last-child(2)>a").Text()

	pageNoAsString = strings.ReplaceAll(pageNoAsString, "Page", "")

	pageNo, err := strconv.Atoi(pageNoAsString)

	if err != nil {
		return 0, errors.New("failed to get page number")
	}

	return pageNo, nil
}

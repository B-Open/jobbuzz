package scraper

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/b-open/jobbuzz/pkg/model"
	"github.com/rs/zerolog/log"
)

func NewBruneidaScraper() BruneidaScraper {
	return BruneidaScraper{wg: sync.WaitGroup{}, jobs: []*model.Job{}, FetchClient: &FetchClient{}}
}

func (s *BruneidaScraper) ScrapeJobs() ([]*model.Job, error) {
	for i := 1; i < 30; i++ {
		s.wg.Add(1)
		url := fmt.Sprintf("https://www.bruneida.com/brunei/jobs/?&page=%d", i)

		go s.scrapeBruneidaJobsListing(url)

	}

	s.wg.Wait()

	return s.jobs, nil
}

func (bruneidaScraper *BruneidaScraper) scrapeBruneidaJobsListing(url string) {
	defer bruneidaScraper.wg.Done()

	links, err := bruneidaScraper.getJobLinks(url)
	if err != nil {
		log.Error().Err(err).Msgf("Fail to scrape url : %s", url)
		return
	}

	for _, link := range links {
		bruneidaScraper.wg.Add(1)

		go bruneidaScraper.scrapeBruneidaJob(link)
	}
}

func (bruneidaScraper *BruneidaScraper) scrapeBruneidaJob(url string) bool {
	defer bruneidaScraper.wg.Done()
	doc, err := bruneidaScraper.FetchClient.GetDocument(url)
	if err != nil {
		log.Error().Err(err).Msgf("Fail to scrape url : %s", url)
		return false
	}

	jobTitle := doc.Find("#title-box-inner div.inline-block.pull-left h1").Text()
	// company := doc.Find("#ad-contact ul li:first-child span.bb b.small").Text()
	salary := doc.Find("#ad-body-inner .opt .opt-dl:nth-child(3) .dd").Text()

	description, err := minifyHtml(doc.Find("#full-description").Text())
	if err != nil {
		log.Error().Err(err).Msgf("Fail to get description : %s", url)
		return false
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

	providerJobId, err := getBruneidaJobId(url)
	if err != nil {
		log.Error().Err(err).Msgf("Fail to get job provider id : %s", url)
		return false
	}

	job := model.Job{
		ProviderJobID: providerJobId,
		Provider:      Bruneida,
		Title:         jobTitle,
		// Company:       company,
		Salary:        salary,
		Location:      location,
		Description:   *description,
	}

	bruneidaScraper.jobs = append(bruneidaScraper.jobs, &job)
	return true

}

func (s *BruneidaScraper) getJobLinks(url string) ([]string, error) {

	links := []string{}

	doc, err := s.FetchClient.GetDocument(url)
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

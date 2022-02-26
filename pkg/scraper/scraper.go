package scraper

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/avast/retry-go"
	"github.com/b-open/jobbuzz/pkg/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
)

const (
	JobCenter = 1
	Bruneida  = 2
)

type (
	Scraper interface {
		ScrapeJobs() ([]*model.Job, error)
	}

	JobCentreScraper struct {
		wg          sync.WaitGroup
		jobs        []*model.Job
		companies   []*model.Company
		FetchClient FetchClienter
	}

	BruneidaScraper struct {
		wg          sync.WaitGroup
		jobs        []*model.Job
		companies   []*model.Company
		FetchClient FetchClienter
	}

	FetchClienter interface {
		GetDocument(url string) (*goquery.Document, error)
	}

	FetchClient struct{}
)

func (c *FetchClient) GetDocument(url string) (*goquery.Document, error) {
	log.Debug().Msgf("Visiting: %s \n", url)
	var doc *goquery.Document

	err := retry.Do(func() error {
		res, err := http.Get(url)
		if err != nil {
			return errors.Wrapf(err, "Error in HTTP GET")
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			errorMessage := fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status)

			return fmt.Errorf("Status is not 200 OK: %s", errorMessage)
		}

		doc, err = goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			return errors.Wrapf(err, "Error creating goquery document")
		}

		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "Error in retry")
	}

	return doc, nil
}

func minifyHtml(s string) (*string, error) {

	m := minify.New()
	m.Add("text/html", &html.Minifier{
		KeepComments:            false,
		KeepWhitespace:          false,
		KeepDocumentTags:        false,
		KeepQuotes:              true,
		KeepEndTags:             false,
		KeepConditionalComments: false,
		KeepDefaultAttrVals:     false,
	})

	minified, err := m.String("text/html", s)
	if err != nil {
		return nil, err
	}

	return &minified, nil
}

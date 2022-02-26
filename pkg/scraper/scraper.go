package scraper

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/avast/retry-go"
	"github.com/b-open/jobbuzz/pkg/model"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
)

const (
	JobCenter = 1
	Bruneida  = 2
)

type scraper struct {
	wg   sync.WaitGroup
	jobs []*model.Job
}

func createScraper() scraper {
	return scraper{wg: sync.WaitGroup{}, jobs: []*model.Job{}}
}

func getDocument(url string) (*goquery.Document, error) {
	fmt.Printf("Visiting: %s \n", url)
	var doc *goquery.Document

	err := retry.Do(func() error {
		res, err := http.Get(url)
		if err != nil {
			return err
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			errorMessage := fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status)

			fmt.Println(errorMessage)

			return errors.New(errorMessage)
		}

		doc, err = goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
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

package scraper

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const (
	JobCenter = 1
	Bruneida  = 2
)

func getDocument(url string) (*goquery.Document, error) {
	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err)

		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		errorMessage := fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status)

		log.Fatal(errorMessage)

		return nil, errors.New(errorMessage)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		return nil, err
	}

	return doc, nil
}

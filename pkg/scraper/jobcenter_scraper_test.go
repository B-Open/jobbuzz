package scraper

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/b-open/jobbuzz/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestScrapeJobs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<ul class="pagination"><li class="page-item"><a href=""><span>Page</span>1</a></li><li class="page-item"></li></ul>`)
		return
	}))
	defer server.Close()

	scraper := NewJobCentreScraper()
	scraper.BaseURL = server.URL

	_, _, err := scraper.ScrapeJobs()

	assert.Nil(t, err, "Error is not nil")
}

func TestScrapeCompanyDetails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open("../../testdata/jobcentre_company_605802.html")
		if err != nil {
			t.Fatal(err)
		}

		reader := bufio.NewReader(f)
		_, err = io.Copy(w, reader)
		if err != nil {
			t.Fatal(err)
		}
		return
	}))
	defer server.Close()

	scraper := NewJobCentreScraper()
	scraper.BaseURL = server.URL

	testID := "605802"
	company := &model.Company{
		Provider:          JobCenter,
		ProviderCompanyID: testID,
		Name:              "CUCKOO INTERNATIONAL (B) SDN BHD",
		Link:              "/web/guest/view-employer/-/employer/605802",
	}
	companies := map[string]*model.Company{}
	companies[testID] = company

	companies = scraper.scrapeCompanyDetails(companies)

	assert.NotNil(t, companies)
	assert.NotNil(t, companies[testID])
	assert.NotEmpty(t, companies[testID].Content, "Company.Content should not be empty")
}

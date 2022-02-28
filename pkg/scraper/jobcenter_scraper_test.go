package scraper

import (
	"bufio"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/b-open/jobbuzz/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFetchClient struct {
	mock.Mock
}

func (c *MockFetchClient) GetDocument(url string) (*goquery.Document, error) {
	args := c.Called(url)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*goquery.Document), args.Error(1)
}

func TestScrapeJobs(t *testing.T) {
	mockClient := &MockFetchClient{}
	pageHtml := `<ul class="pagination"><li class="page-item"><a href=""><span>Page</span>1</a></li><li class="page-item"></li></ul>`
	pageDoc, err := goquery.NewDocumentFromReader(strings.NewReader(pageHtml))
	if err != nil {
		t.Fatal(err)
	}
	mockClient.On("GetDocument", "https://www.jobcentrebrunei.gov.bn/web/guest/search-job?q=&delta=200&start=1").Return(pageDoc, nil)

	scraper := NewJobCentreScraper()
	scraper.FetchClient = mockClient

	_, _, err = scraper.ScrapeJobs()

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

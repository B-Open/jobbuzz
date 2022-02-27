package scraper

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScrapeJobs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<ul class="pagination"><li class="page-item"><a href=""><span>Page</span>1</a></li><li class="page-item"></li></ul>`)
		return
	}))
	defer server.Close()

	scraper := NewJobCentreScraper()
	scraper.BaseURL = server.URL

	_, _, err := scraper.ScrapeJobs()

	assert.Nil(t, err, "Error is not nil")
}

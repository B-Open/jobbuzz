package scraper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetJobcenterJobId(t *testing.T) {
	want := "82563731"
	got, err := getJobcenterJobId("/web/guest/view-job/-/jobs/82563731/security-guard")

	assert.Nil(t, err)
	assert.Equal(t, want, got)
}

func TestGetJobcenterCompanyId(t *testing.T) {
	want := "496468"
	got, err := getJobcenterCompanyId("/web/guest/view-employer/-/employer/496468")

	assert.Nil(t, err)
	assert.Equal(t, want, got)
}

func TestGetBruneidaJobId(t *testing.T) {

	jobId, _ := getBruneidaJobId("https://www.bruneida.com/SHOP-ASSISTANTS-106679")

	if jobId != "106679" {
		t.Errorf("Expected %s but got %s", "expected", "got")
	}
}

func TestGetDocument(t *testing.T) {
	url := "https://github.com"

	client := FetchClient{}
	_, err := client.GetDocument(url)

	if err != nil {
		t.Error(err)
	}

}

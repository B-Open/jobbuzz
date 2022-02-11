package scraper

import "testing"

func TestGetJobcenterJobId(t *testing.T) {
	jobId, _ := getJobcenterJobId("/web/guest/view-job/-/jobs/82563731/security-guard")

	if jobId != "82563731" {
		t.Errorf("Expected %s but got %s", "expected", "got")
	}
}

func TestGetBruneidaJobId(t *testing.T) {

	jobId, _ := getBruneidaJobId("https://www.bruneida.com/SHOP-ASSISTANTS-106679")

	if jobId != "106679" {
		t.Errorf("Expected %s but got %s", "expected", "got")
	}
}

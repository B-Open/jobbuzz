package scraper

import (
	"fmt"

	"github.com/gocolly/colly"
)

func HandleError(collectors []*colly.Collector) {
	for _, collector := range collectors {
		collector.OnError(func(r *colly.Response, err error) {
			fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
		})
	}
}

func HandleRequest(collectors []*colly.Collector) {
	for _, collector := range collectors {
		collector.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL.String())
		})
	}
}

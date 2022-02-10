package main

import (
	"fmt"
	"regexp"
)

func main() {
	// This is just for me to test the scripts
	jobId := getJobId("https://www.bruneida.com/Sales-Girl-Sales-Boy-106590")
	fmt.Println(jobId)
}

func getJobId(s string) string {
	r := regexp.MustCompile(`-(?P<jobId>\d+)`)
	return r.FindStringSubmatch(s)[1]
}

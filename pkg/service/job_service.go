package service

import (
	"log"

	"github.com/b-open/jobbuzz/internal/config"
	"github.com/b-open/jobbuzz/pkg/model"
)

func GetJobs() []model.Job {
	db := config.GetDb()

	var jobs []model.Job

	results := db.Find(&jobs)

	if err := results.Error; err != nil {
		log.Fatal(err)
	}

	return jobs
}

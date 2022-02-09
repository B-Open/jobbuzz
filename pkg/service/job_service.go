package service

import (
	"github.com/b-open/jobbuzz/pkg/model"
)

func (s *Service) GetJobs() ([]model.Job, error) {
	var jobs []model.Job

	results := s.Database.Find(&jobs)

	if err := results.Error; err != nil {
		return nil, err
	}

	if results.RowsAffected < 1 {
		return make([]model.Job, 0), nil
	}

	return jobs, nil
}

func (s *Service) CreateJob(job *model.Job) (*model.Job, error) {
	result := s.Database.Create(&job)

	if err := result.Error; err != nil {
		return nil, err
	}

	return job, nil
}

package service

import (
	"github.com/b-open/jobbuzz/pkg/model"
)

func (s *Service) GetJobs() ([]*model.Job, error) {
	var jobs []*model.Job

	results := s.DB.Find(&jobs)
	if err := results.Error; err != nil {
		return nil, err
	}

	return jobs, nil
}

func (s *Service) CreateJob(job *model.Job) (*model.Job, error) {

	result := s.DB.FirstOrCreate(&job, model.Job{Provider: job.Provider, ProviderJobId: job.ProviderJobId})
	if err := result.Error; err != nil {
		return nil, err
	}

	return job, nil
}

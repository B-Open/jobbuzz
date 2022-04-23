package service

import (
	"github.com/b-open/jobbuzz/pkg/graph/graphmodel"
	"github.com/b-open/jobbuzz/pkg/model"
	"github.com/b-open/jobbuzz/pkg/pagination"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (s *Service) GetJobs(paginationInput graphmodel.PaginationInput) ([]*model.Job, *pagination.PaginationResult, error) {
	var jobs []*model.Job

	results := s.DB.Limit(paginationInput.Size).Offset(pagination.GetOffset(paginationInput)).Find(&jobs)
	if results.Error != nil {
		return nil, nil, results.Error
	}

	var totalCount int64
	countResult := s.DB.Model(&model.Job{}).Where("deleted_at IS NULL").Count(&totalCount)
	if countResult.Error != nil {
		return nil, nil, countResult.Error
	}

	paginationResult := pagination.GetPaginationResult(paginationInput, len(jobs), int(totalCount))

	return jobs, &paginationResult, nil
}

func (s *Service) CreateJobsAndCompanies(jobs []*model.Job, companies map[string]*model.Company) error {
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// create companies
		for _, company := range companies {
			conditions := model.Company{
				Provider:          company.Provider,
				ProviderCompanyID: company.ProviderCompanyID,
			}
			result := tx.FirstOrCreate(&company, conditions)
			if result.Error != nil {
				return errors.Wrapf(result.Error, "Error in Company insert")
			}
		}

		// set company id in jobs
		for _, job := range jobs {
			company, ok := companies[job.Company.ProviderCompanyID]
			if !ok {
				continue
			}
			job.CompanyID = &company.ID
		}

		_, err := s.CreateJobs(tx, jobs)
		if err != nil {
			return errors.Wrapf(err, "Failed to create jobs")
		}

		return nil
	})
	if err != nil {
		return errors.Wrapf(err, "Error in Transaction")
	}

	return nil
}

func (s *Service) CreateJobs(db *gorm.DB, jobs []*model.Job) ([]*model.Job, error) {
	for _, job := range jobs {
		conditions := model.Job{
			Provider:      job.Provider,
			ProviderJobID: job.ProviderJobID,
		}
		result := db.Omit("Company").FirstOrCreate(&job, conditions)
		if result.Error != nil {
			return nil, errors.Wrapf(result.Error, "Error in Job insert")
		}
	}

	return jobs, nil
}

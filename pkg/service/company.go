package service

import (
	"github.com/b-open/jobbuzz/pkg/graph/graphmodel"
	"github.com/b-open/jobbuzz/pkg/model"
)

func (s *Service) GetCompanies(pagination graphmodel.PaginationInput) ([]*model.Company, int64, error) {
	var companies []*model.Company

	results := s.DB.Limit(pagination.Limit).Offset(pagination.Offset).Find(&companies)
	if err := results.Error; err != nil {
		return nil, 0, err
	}

	var totalCount int64
	// TODO: check deleted
	countResult := s.DB.Model(&model.Company{}).Count(&totalCount)
	if countResult.Error != nil {
		return nil, 0, countResult.Error
	}

	return companies, totalCount, nil
}

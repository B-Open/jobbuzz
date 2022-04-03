package service

import (
	"github.com/b-open/jobbuzz/pkg/graph/graphmodel"
	"github.com/b-open/jobbuzz/pkg/model"
)

func (s *Service) GetCompanies(pagination graphmodel.PaginationInput) ([]*model.Company, error) {
	var companies []*model.Company

	results := s.DB.Limit(*pagination.Limit).Offset(*pagination.Offset).Find(&companies)
	if err := results.Error; err != nil {
		return nil, err
	}

	return companies, nil
}

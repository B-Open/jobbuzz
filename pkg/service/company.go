package service

import (
	"github.com/b-open/jobbuzz/pkg/graph/graphmodel"
	"github.com/b-open/jobbuzz/pkg/model"
	"github.com/b-open/jobbuzz/pkg/pagination"
)

func (s *Service) GetCompanies(paginationInput graphmodel.PaginationInput) ([]*model.Company, *pagination.PaginationResult, error) {
	var companies []*model.Company

	results := s.DB.Limit(paginationInput.Size).Offset(pagination.GetOffset(paginationInput)).Find(&companies)
	if err := results.Error; err != nil {
		return nil, nil, err
	}

	var totalCount int64
	countResult := s.DB.Model(&model.Company{}).Where("deleted_at IS NULL").Count(&totalCount)
	if countResult.Error != nil {
		return nil, nil, countResult.Error
	}

	paginationResult := pagination.GetPaginationResult(paginationInput, len(companies), int(totalCount))

	return companies, &paginationResult, nil
}

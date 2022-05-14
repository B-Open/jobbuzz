package service

import (
	"github.com/b-open/jobbuzz/pkg/graph/graphmodel"
	"github.com/b-open/jobbuzz/pkg/model"
	"github.com/b-open/jobbuzz/pkg/pagination"
	"gorm.io/gorm"
)

type Servicer interface {
	GetJobs(pagination graphmodel.PaginationInput) ([]*model.Job, *pagination.PaginationResult, error)
	GetCompanies(pagination graphmodel.PaginationInput) ([]*model.Company, *pagination.PaginationResult, error)
	CreateUser(email string, password string) (token string, err error)
}

type Service struct {
	DB *gorm.DB
}

package service

import (
	"github.com/b-open/jobbuzz/pkg/graph/graphmodel"
	"github.com/b-open/jobbuzz/pkg/model"
	"gorm.io/gorm"
)

type Servicer interface {
	GetJobs(pagination graphmodel.PaginationInput) ([]*model.Job, int64, error)
	GetCompanies(pagination graphmodel.PaginationInput) ([]*model.Company, int64, error)
	CreateUser(email string, password string) (token string, err error)
}

type Service struct {
	DB *gorm.DB
}

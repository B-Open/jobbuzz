package service

import (
	"github.com/b-open/jobbuzz/pkg/graph/graphmodel"
	"github.com/b-open/jobbuzz/pkg/model"
	"gorm.io/gorm"
)

type Servicer interface {
	GetJobs(pagination graphmodel.PaginationInput) ([]*model.Job, error)
}

type Service struct {
	DB *gorm.DB
}

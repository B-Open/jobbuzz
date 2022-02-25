package service

import (
	"github.com/b-open/jobbuzz/pkg/model"
	"gorm.io/gorm"
)

type Servicer interface {
	GetJobs() ([]*model.Job, error)
}

type Service struct {
	DB *gorm.DB
}

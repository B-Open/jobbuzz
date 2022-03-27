package service

import (
	"github.com/b-open/jobbuzz/pkg/model"
	"gorm.io/gorm"
)

type Servicer interface {
	GetJobs() ([]*model.Job, error)
	CreateUser(email string, password string) (token string, err error)
}

type Service struct {
	DB *gorm.DB
}

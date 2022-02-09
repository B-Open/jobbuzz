package service

import "gorm.io/gorm"

type Service struct {
	Database *gorm.DB
}

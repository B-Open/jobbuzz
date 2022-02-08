package model

import "time"

type Job struct {
	ID        uint `gorm:"primaryKey"`
	Title     string
	Company   string
	Salary    string
	Location  string
	CreatedAt time.Time `gorm:"autoUpdateTime"`
	UpdatedAt time.Time `gorm:"autoCreateTime"`
}

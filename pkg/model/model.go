package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

type Job struct {
	BaseModel
	Provider    int    `json:"provider"`
	JobId       string `json:"jobId"`
	Title       string `json:"title"`
	Company     string `json:"company"`
	Salary      string `json:"salary"`
	Location    string `json:"location"`
	Link        string `json:"link"`
	Description string `json:"description"`
}

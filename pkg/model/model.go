package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type Job struct {
	BaseModel
	Title    string `json:"title"`
	Company  string `json:"company"`
	Salary   string `json:"salary"`
	Location string `json:"location"`
}

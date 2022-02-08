package model

import (
	"gorm.io/gorm"
)

type Job struct {
	gorm.Model
	Title    string
	Company  string
	Salary   string
	Location string
}

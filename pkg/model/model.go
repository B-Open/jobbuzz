package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint64         `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

type Job struct {
	BaseModel
	CompanyID     *uint64  `json:"company_id"`
	Company       *Company `json:"company"`
	Provider      int      `gorm:"primarykey;autoIncrement:false" json:"provider"`
	ProviderJobID string   `gorm:"primarykey;autoIncrement:false" json:"providerJobId"`
	Title         string   `json:"title"`
	Salary        string   `json:"salary"`
	Location      string   `json:"location"`
	Link          string   `json:"link"`
	Description   string   `json:"description"`
}

type Company struct {
	BaseModel
	Provider          int    `gorm:"primarykey;autoIncrement:false" json:"provider"`
	ProviderCompanyID string `gorm:"primarykey;autoIncrement:false" json:"providerCompanyId"`
	Name              string `json:"name"`
	Content           string `json:"content"` // TODO: can be refactored to use gorm custom data type
	Link              string `json:"link"`
}

func (c *Company) SetContent(content interface{}) error {
	jsonBytes, err := json.Marshal(&content)
	if err != nil {
		return err
	}
	c.Content = string(jsonBytes)
	return nil
}

func (c *Company) GetContent() (map[string]interface{}, error) {
	var content map[string]interface{}
	err := json.Unmarshal([]byte(c.Content), &content)
	return content, err
}

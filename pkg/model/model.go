package model

import (
	"encoding/json"
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
	CompanyID     uint   `json:"companyId"`
	Provider      int    `json:"provider"`
	ProviderJobID string `json:"providerJobId"`
	Title         string `json:"title"`
	Salary        string `json:"salary"`
	Location      string `json:"location"`
	Link          string `json:"link"`
	Description   string `json:"description"`
}

// TODO: add index
// TODO: add association
type Company struct {
	BaseModel
	ProviderCompanyID string `json:"providerCompanyId"`
	Name              string `json:"name"`
	Content           string `json:"content"`
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
	// TODO: return strictly typed struct
	var content map[string]interface{}
	err := json.Unmarshal([]byte(c.Content), &content)
	return content, err
}

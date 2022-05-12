package models

import (
	"fmt"
	"github.com/lib/pq"
	"next/utils"
)

type Metadata struct {
	BaseModel
	Name   string         `json:"name"`
	Values pq.StringArray `json:"values" gorm:"type:text;null"`
}

func (m *Metadata) Create() error {
	fmt.Printf("Data %v\n", m)
	if err := db.Create(&m).Error; err != nil {
		utils.LogError("Create metadata", err)
		return err
	}
	return nil
}

func (m *Metadata) Update(id string) error {
	if err := db.Model(&Metadata{}).Where("id = ?", id).Updates(&m).Error; err != nil {
		utils.LogError("Update metadata", err)
		return err
	}
	return nil
}

func (m *Metadata) Delete(id string) error {
	if err := db.Where("id = ?", id).Delete(&m).Error; err != nil {
		utils.LogError("Delete metadata", err)
		return err
	}
	return nil
}

func (m *Metadata) GetAll() (metadata []Metadata, err error) {
	if err = db.Debug().Find(&metadata).Error; err != nil {
		utils.LogError("Get all metadata", err)
		return
	}
	return
}

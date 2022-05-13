package models

import (
	uuid "github.com/satori/go.uuid"
)

type Stock struct {
	BaseModel
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid"`
	Quantity  int64     `json:"quantity" gorm:"default:0"`
	Variant   string    `json:"variant"`
}

func (s *Stock) Create() error {
	return db.Create(s).Error
}

func (s *Stock) Update() error {
	return db.Model(s).Where("product_id = ? AND variant = ?", s.ProductID, s.Variant).Update("quantity", s.Quantity).Error
}

func (s *Stock) Get(productId string, variant string) error {
	return db.Debug().Where("product_id = ? AND variant = ?", productId, variant).First(s).Error
}

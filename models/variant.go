package models

import (
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

type Variant struct {
	BaseModel
	Product   Product        `json:"-" gorm:"foreignkey:ProductID;constraint:OnDelete:CASCADE"`
	ProductID uuid.UUID      `json:"-"`
	Key       string         `json:"key"`
	Values    pq.StringArray `json:"values" gorm:"type:text"`
}

//func (variant *Variant) Create(product_id string, varis []Variant) ([]Variant, error) {
//	for index, _ := range varis {
//		varis[index].ProductID = uuid.FromStringOrNil(product_id)
//	}
//	err := db.Create(varis).Error
//	if err != nil {
//		return nil, err
//	}
//
//	var stocks []Stock
//
//	productUUID, err := uuid.FromString(product_id)
//	if err != nil {
//		return nil, err
//	}
//
//	for _, combinations := range Combination(varis) {
//		stocks = append(stocks, Stock{
//			ProductID: productUUID,
//			Variant:   combinations,
//		})
//	}
//
//	stock.
//	return varis, nil
//}

func (variant *Variant) Update(id string) error {
	return db.Model(&Variant{}).Where("id = ?", id).Updates(&variant).Error
}

func (variant *Variant) Delete(id string) error {
	return db.Where("id = ?", id).Delete(&Variant{}).Error
}

func (variant *Variant) Get(product_id string) (variants []Variant, err error) {
	err = db.Where("product_id = ?", product_id).Find(&variants).Error
	return
}

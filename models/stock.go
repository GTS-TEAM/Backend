package models

import (
	"errors"
	uuid "github.com/satori/go.uuid"
	"net/url"
)

type Stock struct {
	BaseModel
	Quantity  int64     `json:"quantity"`
	Product   Product   `json:"product"`
	ProductID uuid.UUID `json:"product_id"`
	Metadata  JSONB     `json:"metadata" gorm:"type:jsonb;null"`
}

func (stock *Stock) Create() error {
	return db.Create(stock).Error
}

func (stock *Stock) Update(id string) error {
	return db.Model(&Stock{}).Where("id = ?", id).Updates(&stock).Error
}

func (stock *Stock) Delete(id string) error {
	return db.Where("id = ?", id).Delete(&Stock{}).Error
}

func (stock *Stock) Get(q url.Values) (int64, error) {

	var query string
	productId := q.Get("product_id")
	if productId == "" {
		return 0, errors.New("product_id is required")
	}
	q.Del("product_id")
	// attrs->>'name' = 'hello';
	for k, v := range q {
		query += " AND metadata->>'" + k + "' = '" + v[0] + "'"
	}

	err := db.Model(&Stock{}).Where("product_id = ?"+query, productId).First(&stock).Error

	if err != nil {
		return 0, err
	}
	return stock.Quantity, nil
}

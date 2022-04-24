package models

import (
	"fmt"
	"github.com/lib/pq"
)

type Product struct {
	BaseModel
	Name string `json:"name"`
	Price float64 `json:"price"`
	Description string `json:"description"`
	Fields JSONB `gorm:"type:jsonb;not null" json:"fields"`
	Quantity int `json:"quantity"`
	Images pq.StringArray `gorm:"type:text" json:"images"`
	//CategoryId uuid.UUID `json:"category_id"`
	//Category Category `json:"category"`
}

func (p *Product) Create(dto Product) error {

	fmt.Println("Create Product: ", dto)

	if  err := db.Create(&dto).Error; err != nil {
		fmt.Printf("Error creating product: %v", err)
		return err
	}
	return nil
}

func (p Product) GetAll(dto Product) (products []Product) {
	fmt.Println("Gets Product: ", dto)
	db.Find(&products)
	return
}

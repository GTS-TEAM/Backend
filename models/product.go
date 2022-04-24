package models

import (
	"fmt"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

type Product struct {
	BaseModel
	Name string `json:"name"`
	Price float64 `json:"price"`
	Description string `json:"description"`
	Fields JSONB `gorm:"type:jsonb;not null" json:"fields"`
	Quantity int `json:"quantity"`
	Images pq.StringArray `gorm:"type:text" json:"images"`
	Categories  []*Category   `json:"categories,omitempty" gorm:"many2many:products_categories"`
	CategoriesId []uuid.UUID `gorm:"-" json:"-"`
	//Category Category `json:"category"`
}

func (p *Product) Create(dto Product) error {

	var categories []*Category

	db.Model(Category{}).Where("id IN (?)", dto.CategoriesId).Find(&categories)
	dto.Categories = categories

	fmt.Println("Create Product: ", dto)

	if  err := db.Omit("Categories.*").Create(&dto).Error; err != nil {
		fmt.Printf("Error creating product: %v", err)
		return err
	}
	return nil
}

func (p Product) GetAll(category string) (products []Product) {
	err := db.Where("categories.id = ?", category).
		Joins("JOIN products_categories ON products_categories.product_id = products.id").
		Joins("JOIN categories ON categories.id = products_categories.category_id").
		Find(&products).Error

	if err != nil {
		fmt.Printf("Error getting products: %v", err)
	}

	return
}

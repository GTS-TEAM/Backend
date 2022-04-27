package models

import (
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm/clause"
	"next/utils"
)

type ProductResponse struct {
	BaseModel
	Name        string         `json:"name"`
	Price       float64        `json:"price"`
	Description string         `json:"description"`
	Metadata    JSONB          `json:"metadata"`
	Quantity    int            `json:"quantity"`
	Images      pq.StringArray `json:"images"`
	Categories  []*Category    `json:"categories,omitempty"`
	User        *User          `json:"creator"`
}

type ProductsResponse struct {
	Products []ProductResponse `json:"products"`
	Total    int64             `json:"total"`
}

type Product struct {
	BaseModel
	Name         string         `json:"name"`
	Price        float64        `json:"price"`
	Description  string         `json:"description"`
	Metadata     JSONB          `gorm:"type:jsonb;not null" json:"metadata"`
	Quantity     int            `json:"quantity"`
	Images       pq.StringArray `gorm:"type:text" json:"images"`
	Categories   []*Category    `json:"categories,omitempty" gorm:"many2many:products_categories"`
	CategoriesId []uuid.UUID    `gorm:"-" json:"categories_id"`
	UserID       uuid.UUID      `gorm:"type:uuid" json:"-"`
	User         *User          `json:"creator" gorm:"foreignkey:UserID"`
}

func (p *Product) Create(userId string, dto *Product) error {

	if err := db.Model(Category{}).Where("id IN (?)", dto.CategoriesId).Find(&dto.Categories).Error; err != nil {
		return err
	}

	if err := db.Model(User{}).Where("id = ?", userId).First(&User{}).Error; err != nil {
		return err
	}

	if err := db.Omit("Categories.*").Create(&dto).Error; err != nil {
		utils.LogError(err)
		return err
	}
	return nil
}

func (p Product) GetAll(category string, paging Pagination) (data ProductsResponse, err error) {

	var products []Product

	if category == "" {
		err = db.Offset(paging.Page).Limit(paging.Limit).Order(paging.Sort).Find(&products).Count(&data.Total).Error
		if err != nil {
			utils.LogError(err)
		}
		return
	} else {
		db.Preload(clause.Associations).
			Where("categories.id = ? and categories.deleted_at IS NULL ", category).
			Joins("LEFT JOIN products_categories ON products_categories.product_id = products.id").
			Joins("LEFT JOIN categories ON categories.id = products_categories.category_id").
			Where("products.deleted_at IS NULL").
			Offset(paging.Page).
			Limit(paging.Limit).
			Order("products." + paging.Sort).Find(&products).Count(&data.Total)
	}

	err = utils.BindStruct(&products, &data.Products)
	if err != nil {
		return ProductsResponse{}, err
	}
	return
}

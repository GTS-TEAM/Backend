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
	Rating      float64        `json:"rating"`
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
	Stock        int            `json:"stock"`
	Specific     JSONB          `json:"specific" gorm:"type:jsonb;null"`
}

func (p *Product) Create(userId string, dto *Product) error {

	if err := db.Model(Category{}).Where("id IN (?)", dto.CategoriesId).Find(&dto.Categories).Error; err != nil {
		return err
	}

	if err := db.Model(User{}).Where("id = ?", userId).First(&User{}).Error; err != nil {
		return err
	}
	dto.UserID = uuid.FromStringOrNil(userId)

	if err := db.Omit("Categories.*").Create(&dto).Error; err != nil {
		utils.LogError("Product Create", err)
		return err
	}
	return nil
}

func (p *Product) GetAll(category string, paging Pagination) (data ProductsResponse, err error) {

	var products []Product

	if category == "" {
		err = db.Preload(clause.Associations).Offset(paging.Page).Limit(paging.Limit).Order(paging.Sort).Find(&products).Count(&data.Total).Error
		if err != nil {
			utils.LogError("Product GetAll", err)
		}

		err = utils.BindStruct(&products, &data.Products)
		if err != nil {
			utils.LogError("Product GetAll", err)
		}
	} else {
		db.Preload(clause.Associations).
			Where("categories.id = ? and categories.deleted_at IS NULL ", category).
			Joins("LEFT JOIN products_categories ON products_categories.product_id = products.id").
			Joins("LEFT JOIN categories ON categories.id = products_categories.category_id").
			Where("products.deleted_at IS NULL").
			Offset(paging.Page).Limit(paging.Limit).
			Order("products." + paging.Sort).Find(&products).Count(&data.Total)
	}

	err = utils.BindStruct(&products, &data.Products)

	for index, _ := range data.Products {
		err = db.Model(&Review{}).Select("AVG(rating) as rating").Where("product_id = ?", data.Products[index].ID).Scan(&data.Products[index].Rating).Error
		if err != nil {
			data.Products[index].Rating = float64(0)
			err = nil
		}
	}

	if err != nil {
		return ProductsResponse{}, err
	}
	return
}

func (p *Product) GetByID(id string) (data ProductResponse, err error) {
	var product Product

	err = db.Model(&Product{}).Preload(clause.Associations).Where("id = ?", id).First(&product).Error

	if err != nil {
		utils.LogError("Product GetByID", err)
		return
	}
	err = utils.BindStruct(&product, &data)
	if err != nil {
		utils.LogError("Product GetByID", err)
		return
	}
	err = db.Model(&Review{}).Select("AVG(rating) as rating").Where("product_id = ?", data.ID).Scan(&data.Rating).Error
	if err != nil {
		data.Rating = float64(0)
		err = nil
	}
	return
}

func (p *Product) Update(id string, dto *Product) (err error) {
	err = db.Preload(clause.Associations).Where("id = ?", id).First(&p).Error
	if err != nil {
		utils.LogError("Product Update", err)
		return
	}
	err = utils.BindStruct(&dto, &p)
	if err != nil {
		utils.LogError("Product Update", err)
		return
	}

	err = db.Where("id = ?", id).Updates(&p).Error

	if err != nil {
		utils.LogError("Product Update", err)
		return
	}
	return nil
}

func (p *Product) Delete(id string) (err error) {
	err = db.Where("id = ?", id).Delete(&Product{}).Error
	if err != nil {
		utils.LogError("Product Delete", err)
		return
	}
	return nil
}

package models

import (
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm/clause"
	"next/utils"
	"strings"
)

type ProductResponse struct {
	BaseModel
	Name        string                   `json:"name"`
	Price       float64                  `json:"price"`
	Description string                   `json:"description"`
	Metadata    []map[string]interface{} `json:"metadata"`
	Stock       int                      `json:"stock"`
	Images      pq.StringArray           `json:"images"`
	Categories  []*Category              `json:"categories,omitempty"`
	User        *User                    `json:"creator"`
	Rating      float64                  `json:"rating"`
}

type ProductsResponse struct {
	Products []ProductResponse `json:"products"`
	Total    int64             `json:"total"`
}

type Product struct {
	BaseModel
	Name         string                   `json:"name"`
	Price        float64                  `json:"price"`
	Description  string                   `json:"description"`
	Images       pq.StringArray           `gorm:"type:text" json:"images"`
	Categories   []*Category              `json:"categories,omitempty" gorm:"many2many:products_categories"`
	CategoriesId []uuid.UUID              `gorm:"-" json:"categories_id"`
	UserID       uuid.UUID                `gorm:"type:uuid" json:"-"`
	User         *User                    `json:"creator" gorm:"foreignkey:UserID"`
	Specific     JSONB                    `json:"specific" gorm:"type:jsonb;null"`
	Metadata     []map[string]interface{} `json:"metadata" gorm:"-"`
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
		err = db.Preload(clause.Associations).
			Offset(paging.Page).
			Limit(paging.Limit).
			Order(paging.Sort).
			Find(&products).
			Count(&data.Total).Error
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

	// TODO: Get quantity

	for i, _ := range products {
		var metadata []map[string]interface{}

		db.Raw("select key, array_agg(distinct value) from stocks join lateral (select * from jsonb_each_text(metadata)) j on true group by key").
			Where("product_id = ?", products[i].ID).
			Scan(&metadata)

		for _, v := range metadata {
			ar := v["array_agg"].(string)
			ar = ar[1 : len(ar)-1]

			var values []map[string]interface{}

			for _, va := range strings.Split(ar, ",") {
				values = append(values, map[string]interface{}{
					"value":    va,
					"quantity": 0,
				})
			}

			m := map[string]interface{}{
				"key":    v["key"],
				"values": values,
			}

			/*
				[
					{
						"key": "size",
						"values": [{"key":"S","quantity":10}, {"key":"M","quantity":20}]
					}

				]

			*/

			for _, value := range m["values"].([]map[string]interface{}) {
				/*
					{
								"key": "size",
								"values": [{"key":"S","quantity":10}, {"key":"M","quantity":20}]
					}
				*/
				var count int64
				db.Debug().Model(&Stock{}).Select("sum(quantity)").Where("product_id = ? AND metadata->>'"+m["key"].(string)+"' = ?", products[i].ID, value["value"].(string)).Scan(&count)
				value["quantity"] = count
			}

			products[i].Metadata = append(products[i].Metadata, m)
		}
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

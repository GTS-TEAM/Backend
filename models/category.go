package models

import (
	uuid "github.com/satori/go.uuid"
	"next/utils"
)

type Category struct {
	BaseModel
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Children    []Category `json:"children,omitempty" gorm:"foreignkey:ParentID"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	Products    []*Product `json:"-" gorm:"many2many:products_categories;constraint:OnDelete:CASCADE"`
}

func (c *Category) TableName() string {
	return "categories"
}

func (c *Category) GetAll(paging Pagination) (result []Category, err error) {

	if err = db.
		Preload("Children").
		Where("parent_id is null").
		Limit(paging.Limit).
		Offset(paging.Page).
		Order("categories." + paging.Sort).
		Find(&result).Error; err != nil {
		return result, err
	}

	return result, nil
}

func (c *Category) GetById(id int) (*Category, error) {
	var category Category
	err := db.First(&category, id).Error
	return &category, err
}

func (c *Category) Create() error {
	return db.Debug().Create(&c).Error
}

func (c *Category) GetCountProductOfCategory(id string) (count int64, err error) {
	err = db.Table("products_categories").Where("category_id = ?", id).Count(&count).Error
	if err != nil {
		utils.LogError("GetCountProductOfCategory", err)
	}
	return
}

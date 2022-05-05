package models

import uuid "github.com/satori/go.uuid"

type CategoryResponse struct {
	Categories []Category `json:"categories"`
	Children   []Category `json:"children"`
	Total      int64      `json:"total"`
}

type Category struct {
	BaseModel
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Children    []Category `json:"children,omitempty" gorm:"foreignkey:ParentID"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	Products    []*Product `json:"-" gorm:"many2many:products_categories"`
}

func (c *Category) TableName() string {
	return "categories"
}

func (c *Category) GetAll(paging Pagination) (result CategoryResponse, err error) {
	//if err = db.Offset(paging.Page).
	//	Limit(paging.Limit).
	//	Order("categories." + paging.Sort).
	//	Find(&result.Categories).Count(&result.Total).Error; err != nil {
	//	return result, err
	//}

	// get parent categories and children categories
	if err = db.
		Preload("Children").
		Where("parent_id is null").
		Limit(paging.Limit).
		Offset(paging.Page).
		Order("categories." + paging.Sort).
		Find(&result.Categories).Error; err != nil {
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

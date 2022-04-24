package models


type Category struct {
	BaseModel
	Name string `json:"name"`
	Description string 	`json:"description"`
	Products []*Product `json:"-" gorm:"many2many:products_categories"`
}

func (c *Category) TableName() string {
	return "categories"
}

func (c *Category) GetAll() ([]Category, error) {
	var categories []Category
	if err := db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (c *Category) GetById(id int) (*Category, error) {
	var category Category
	err := db.First(&category, id).Error
	return &category, err
}

func (c *Category) Create(cate Category) error {
	return db.Create(&cate).Error
}

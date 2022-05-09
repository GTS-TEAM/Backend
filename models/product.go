package models

import (
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm/clause"
	"next/utils"
)

type VariantProductResponse struct {
	Key    string                 `json:"key"`
	Values []ValueVariantResponse `json:"values"`
}

type ValueVariantResponse struct {
	Value    string `json:"value"`
	Quantity int64  `json:"quantity"`
}

type ProductResponse struct {
	BaseModel
	Name        string                   `json:"name"`
	Price       float64                  `json:"price"`
	Description string                   `json:"description"`
	Variants    []VariantProductResponse `json:"variants,omitempty" gorm:"-"`
	Images      interface{}              `json:"images" gorm:"type:text"`
	Specific    interface{}              `json:"specific,omitempty"`
	Categories  string                   `json:"category,omitempty"`
	User        *User                    `json:"creator" gorm:"-"`
	Rating      float64                  `json:"rating"`
}

type ProductsResponse struct {
	Products []Product `json:"products"`
	Total    int64     `json:"total"`
}

type Product struct {
	BaseModel
	Name         string         `json:"name"`
	Price        float64        `json:"price"`
	Description  string         `json:"description"`
	Images       pq.StringArray `gorm:"type:text" json:"images"`
	Categories   []*Category    `json:"categories,omitempty" gorm:"many2many:products_categories"`
	CategoriesId []uuid.UUID    `gorm:"-" json:"categories_id,omitempty"`
	UserID       uuid.UUID      `gorm:"type:uuid" json:"-"`
	User         *User          `json:"creator" gorm:"foreignkey:UserID"`
	Specific     JSONB          `json:"specific" gorm:"type:jsonb;null"`
	Variants     []Variant      `json:"variants,omitempty" gorm:"migration"`
	Rating       float64        `json:"rating" gorm:"migration"`
	Category     string         `json:"category" gorm:"migration"`
}

func Combination(arr []Variant) []string {
	if len(arr) == 1 {
		return arr[0].Values
	}

	var result []string
	var allCasesOfRest = Combination(arr[1:])
	for i := 0; i < len(allCasesOfRest); i++ {
		for j := 0; j < len(arr[0].Values); j++ {
			result = append(result, arr[0].Values[j]+allCasesOfRest[i])
		}
	}
	return result
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

	var variants []Variant
	var stock []Stock

	// TODO : Use S in Solid

	for _, variant := range dto.Variants {
		variants = append(variants, Variant{
			Key:       variant.Key,
			Values:    variant.Values,
			ProductID: dto.ID,
		})
	}

	// TODO: Use S

	for _, combinations := range Combination(dto.Variants) {
		stock = append(stock, Stock{
			ProductID: dto.ID,
			Variant:   combinations,
		})
	}

	if err := db.Create(&variants).Error; err != nil {
		panic(err)
	}

	if err := db.Create(&stock).Error; err != nil {
		panic(err)
	}

	return nil
}

func (p *Product) GetAll(category string, paging Pagination) (data ProductsResponse, result []map[string]interface{}, err error) {

	queryBuilder := db.Debug().Model(&Product{}).
		Select("products.*,User,categories.name as category,avg(reviews.rating) as rating").
		Group("products.id,\"User\".\"id\",categories.name").
		Joins("User").
		Joins("LEFT JOIN reviews ON products.id = reviews.product_id").
		Joins("LEFT JOIN products_categories ON products_categories.product_id = products.id").
		Joins("LEFT JOIN categories ON categories.id = products_categories.category_id")
	if category != "" {
		queryBuilder = queryBuilder.Where("categories.id = ? and categories.deleted_at IS NULL ", category)
	}
	err = queryBuilder.Offset(paging.Page).
		Limit(paging.Limit).
		Order("products.created_at desc").
		Find(&data.Products).
		Count(&data.Total).
		Error

	if err != nil {
		return ProductsResponse{}, nil, err
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

	var variants []Variant

	// Trả về list variant của product
	// [{Key: 'color', Values: ['red', 'blue']}, {Key: 'size', Values: ['S', 'M', 'L']}]
	err = db.Model(&Variant{}).Where("product_id = ?", id).Find(&variants).Error
	if err != nil {
		utils.LogError("Product GetByID", err)
	}

	for _, variant := range variants {
		variantRes := VariantProductResponse{
			Key: variant.Key,
		}
		for _, value := range variant.Values {
			stock := Stock{}
			db.Debug().Model(&Stock{}).Select("sum (quantity)").Where("variant LIKE ?", "%"+value+"%").Scan(&stock.Quantity)

			variantRes.Values = append(variantRes.Values, ValueVariantResponse{
				Value:    value,
				Quantity: stock.Quantity,
			})
		}
		data.Variants = append(data.Variants, variantRes)
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

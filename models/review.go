package models

import (
	"errors"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

type Review struct {
	BaseModel
	UserId    uuid.UUID      `json:"user_id"`
	User      *User          `json:"commenter"`
	Product   *Product       `json:"-"`
	ProductId uuid.UUID      `json:"product_id" gorm:"type:uuid"`
	Rating    int64          `json:"rating" gorm:"default:0"`
	Comment   string         `json:"comment"`
	Images    pq.StringArray `gorm:"type:text" json:"images"`
}

func (r *Review) Create(userId string) error {
	r.UserId, _ = uuid.FromString(userId)
	if err := db.Create(&r).Error; err != nil {
		return err
	}

	return nil
}

func (r *Review) GetReviewOfProduct(productId string) ([]Review, error) {
	var reviews []Review

	if err := db.Preload("User").Where("product_id = ?", productId).Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *Review) DeleteReview(userId string, reviewId string) error {
	var review Review
	if err := db.Where("id = ?", reviewId).First(&review).Error; err != nil {
		return err
	}

	if review.UserId.String() != userId {
		return errors.New("You can't delete this review")
	}

	if err := db.Delete(&review).Error; err != nil {
		return err
	}

	return nil
}

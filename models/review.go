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

type ReviewStatistics struct {
	Count      int64 `json:"count"`
	Rating     int64 `json:"rating"`
	Percentage int64 `json:"percentage"`
}

type ReviewResponse struct {
	Reviews []Review           `json:"reviews"`
	Stats   []ReviewStatistics `json:"stats"`
	Count   int64              `json:"count"`
}

func (r *Review) Create(userId string) error {
	r.UserId, _ = uuid.FromString(userId)
	if err := db.Create(&r).Error; err != nil {
		return err
	}

	return nil
}

func (r *Review) GetReviewOfProduct(productId string, pagination Pagination) (ReviewResponse, error) {
	var reviews []Review
	var statistic []ReviewStatistics
	var count int64

	if err := db.Debug().Model(&Review{}).Preload("User").
		Where("product_id = ?", productId).
		Count(&count).
		Limit(pagination.Limit).
		Offset(pagination.Page).
		Order(pagination.Sort).
		Find(&reviews).Error; err != nil {
		return ReviewResponse{}, err
	}

	db.Model(&Review{}).
		Select("reviews.rating, COUNT(*) AS count").
		Group("rating").
		Where("product_id = ?", productId).
		Find(&statistic)
	// calculate percentage of rating
	for i := range statistic {
		statistic[i].Percentage = int64(float64(statistic[i].Count) / float64(count) * 100)
	}
	return ReviewResponse{
		Reviews: reviews,
		Stats:   statistic,
		Count:   count,
	}, nil
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

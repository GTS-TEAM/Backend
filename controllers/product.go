package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"next/dtos"
	"next/models"
	"next/utils"
)

type ProductController struct {
}

func (p *ProductController) Create(c *gin.Context) {
	userId := getUserID(c)

	product := models.Product{}

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := product.Create(userId, &product)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, dtos.Response{Message: "Product created successfully", Data: product})
}

func (p *ProductController) GetProductsByCategory(c *gin.Context) {
	paging := models.GeneratePaginationFromRequest(c)
	product := models.Product{}
	category := c.Query("category_id")
	filter := utils.GetProductFilter(c)

	fmt.Printf("Filter %+v\n", filter)

	data, err := product.GetAll(category, filter, paging)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, dtos.Response{
		Message: "Success",
		Data:    data,
	})
}

func (p *ProductController) GetProductById(c *gin.Context) {
	product := models.Product{}
	id := c.Param("id")
	prod, err := product.GetByID(id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, dtos.Response{
		Message: "Success",
		Data:    prod,
	})
}

func (p *ProductController) Update(c *gin.Context) {
	product := models.Product{}
	id := c.Param("id")
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := product.Update(id, &product)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, dtos.Response{Message: "Product updated successfully", Data: product})
}

/**
 * Review
 */

func (p *ProductController) CreateReviews(c *gin.Context) {
	review := models.Review{}

	userId := getUserID(c)

	err := c.ShouldBindJSON(&review)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = review.Create(userId)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, dtos.Response{
		Message: "Success",
		Data:    review,
	})
}

func (p *ProductController) GetReviews(c *gin.Context) {
	productId := c.Param("id")
	review := models.Review{}

	reviews, err := review.GetReviewOfProduct(productId)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, dtos.Response{
		Message: "Success",
		Data:    reviews,
	})
}

func (p *ProductController) Delete(c *gin.Context) {
	product := models.Product{}
	id := c.Param("id")
	err := product.Delete(id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, dtos.Response{Message: "Product deleted successfully"})
}

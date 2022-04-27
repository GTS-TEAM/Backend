package controllers

import (
	"github.com/gin-gonic/gin"
	"next/dtos"
	"next/models"
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
	c.ShouldBindJSON(&product)
	category := c.Param("id")

	products, err := product.GetAll(category, paging)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, dtos.Response{
		Message: "Success",
		Data:    products,
	})
}

package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"next/dtos"
	"next/models"
)

type ProductController struct {
}

func (p *ProductController) Create(c *gin.Context)  {

	product := models.Product{}


	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	product.Create(product)
}

func (p *ProductController) GetAll(c *gin.Context)  {

	product := models.Product{}
	c.ShouldBindJSON(&product)
	category := c.Param("id")
	fmt.Println("category", category)
	products := product.GetAll(category)

	c.JSON(200, dtos.Response{
		Message: "Success",
		Data:    products,
	})
}

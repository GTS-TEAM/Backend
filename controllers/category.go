package controllers

import (
	"github.com/gin-gonic/gin"
	"next/dtos"
	"next/models"
)

type CategoryController struct {
}

func (ca CategoryController) Create(c *gin.Context) {

	category := models.Category{}
	err := c.ShouldBindJSON(&category)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = category.Create()
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, dtos.Response{
		Message: "Category created successfully",
		Data:    category,
	})
}

func (ca CategoryController) GetAll(c *gin.Context) {

	category := models.Category{}
	paging := models.GeneratePaginationFromRequest(c)

	categories, err := category.GetAll(paging)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, dtos.Response{
		Message: "Categories retrieved successfully",
		Data:    categories,
	})
}

package controllers

import (
	"github.com/gin-gonic/gin"
	"next/dtos"
	"next/models"
)

type VariantController struct {
}

/*
	Function: CreateStock
	Description: Creates a new stock
	Input: stock: Variant
	Output: Variant
	{
		"metadata":{
			"color": "red",
			"size": "small",
		}
		"quantity": 10,
		"product_id": "f8f8f8f8-f8f8-f8f8-f8f8-f8f8f8f8f8f8",
	}
**/
//func (s *VariantController) Create(c *gin.Context) {
//
//	variant := models.Variant{}
//
//	productId, exists := c.GetQuery("product_id")
//	if !exists {
//		c.JSON(400, gin.H{
//			"message": "product_id is required",
//		})
//	}
//	var variants []models.Variant
//	if err := c.ShouldBindJSON(&variants); err != nil {
//		c.JSON(400, gin.H{"error": err.Error()})
//		return
//	}
//
//	variantsRes, err := variant.Create(productId, variants)
//	if err != nil {
//		c.JSON(400, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(201, dtos.Response{
//		Message: "Variant created successfully",
//		Data:    variantsRes,
//	})
//}

func (s *VariantController) Get(c *gin.Context) {

	variant := models.Variant{}
	productId, _ := c.GetQuery("product_id")

	variants, err := variant.Get(productId)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, dtos.Response{
		Message: "Variant found successfully",
		Data:    variants,
	})
}

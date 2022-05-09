package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"next/dtos"
	"next/models"
)

type StockController struct {
}

func (s *StockController) Get(c *gin.Context) {

	stock := &models.Stock{}

	productId, ok := c.GetQuery("product_id")
	if ok != true {
		c.JSON(400, gin.H{
			"error": "product_id is required",
		})
	}
	fmt.Printf("%v", productId)

	variant, ok := c.GetQuery("variant")
	if ok != true {
		c.JSON(400, gin.H{
			"error": "variant is required",
		})
		return
	}

	if err := stock.Get(productId, variant); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, dtos.Response{
		Data:    stock,
		Message: "success",
	})
}

func (s *StockController) Update(c *gin.Context) {

	stock := &models.Stock{}

	if err := c.ShouldBindJSON(stock); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
	}

	if err := stock.Update(); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(200, dtos.Response{
		Data:    stock,
		Message: "Update stock success",
	})
}

package controllers

import (
	"github.com/gin-gonic/gin"
	"next/dtos"
	"next/models"
)

type StockController struct {
}

/*
	Function: CreateStock
	Description: Creates a new stock
	Input: stock: Stock
	Output: Stock
	{
		"metadata":{
			"color": "red",
			"size": "small",
		}
		"quantity": 10,
		"product_id": "f8f8f8f8-f8f8-f8f8-f8f8-f8f8f8f8f8f8",
	}
**/
func (s *StockController) Create(c *gin.Context) {

	stock := models.Stock{}
	if err := c.ShouldBindJSON(&stock); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := stock.Create()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, dtos.Response{
		Message: "Stock created successfully",
		Data:    stock,
	})
}

func (s *StockController) Get(c *gin.Context) {

	stock := models.Stock{}
	q := c.Request.URL.Query()

	quantity, err := stock.Get(q)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, dtos.Response{
		Message: "Stock found successfully",
		Data: struct {
			Quantity int64 `json:"quantity"`
		}{
			Quantity: quantity,
		},
	})
}

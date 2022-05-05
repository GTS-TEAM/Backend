package controllers

import (
	"github.com/gin-gonic/gin"
	"next/dtos"
	"next/models"
)

type MetadataController struct {
}

func (m *MetadataController) Create(c *gin.Context) {
	metadata := models.Metadata{}

	if err := c.ShouldBindJSON(&metadata); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := metadata.Create(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}

	c.JSON(200, dtos.Response{
		Message: "success",
		Data:    metadata,
	})
}

func (m *MetadataController) Update(c *gin.Context) {
	metadata := models.Metadata{}
	id := c.Param("id")

	if err := c.ShouldBindJSON(&metadata); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := metadata.Update(id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}

	c.JSON(200, dtos.Response{
		Message: "success",
		Data:    nil,
	})
}

func (m *MetadataController) GetAll(c *gin.Context) {
	metadata := models.Metadata{}

	data, err := metadata.GetAll()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}

	c.JSON(200, dtos.Response{
		Message: "success",
		Data:    data,
	})
}

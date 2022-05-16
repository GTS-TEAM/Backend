package controllers

import (
	"github.com/gin-gonic/gin"
	"next/models"
)

type UserController struct {
}

func getUserID(c *gin.Context) (userID string) {
	//MustGet returns the value for the given key if it exists, otherwise it panics.
	// get x-user-id from header
	return c.Request.Header.Get("x-user-id")
}

func (u *UserController) Get(c *gin.Context) {
	userId := getUserID(c)
	user := models.User{}

	name, err := user.GetName(userId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"name": name})
}

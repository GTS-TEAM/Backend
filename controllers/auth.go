package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"next/dtos"
	"next/models"
)

type AuthController struct {
}

func (auth *AuthController) Login(c *gin.Context) {
	var loginForm dtos.LoginForm

	if validationErr := c.ShouldBindJSON(&loginForm); validationErr != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": validationErr})
		return
	}

	user := models.User{}

	res, err := user.Login(loginForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, dtos.Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dtos.Response{Message: "success", Data: res})
}

func (auth *AuthController) Register(c *gin.Context) {
	RegisterForm := dtos.RegisterForm{}

	if validationErr := c.ShouldBindJSON(&RegisterForm); validationErr != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": validationErr})
		return
	}

	user := models.User{}

	err := user.Register(RegisterForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, dtos.Response{Message: err.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusOK, dtos.Response{Message: "Register success", Data: nil})
}

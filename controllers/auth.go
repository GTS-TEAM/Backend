package controllers

import (
	"fmt"
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

func (auth *AuthController) RefreshToken(c *gin.Context) {
	RefreshTokenForm := dtos.RefreshTokenForm{}
	if validationErr := c.ShouldBindJSON(&RefreshTokenForm); validationErr != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": validationErr})
		return
	}

	token := models.Token{}

	res, err := token.ValidateTokenRefreshToken(RefreshTokenForm.RefreshToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, dtos.Response{Message: err.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusOK, dtos.Response{Message: "success", Data: res})
}

func (auth *AuthController) Authorize(c *gin.Context) {
	fmt.Println("Header: ", c.Request.Header)
	t := &models.Token{}
	token := t.ExtractToken(c.Request)

	jwtClaims, err := t.VerifyToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Header("x-user-id", jwtClaims.UserID.String())
	c.Header("x-role", jwtClaims.Role)
	c.JSON(http.StatusOK, dtos.Response{Message: "success", Data: nil})
}

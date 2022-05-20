package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"next/dtos"
	"next/middlewares"
	"next/models"
	"next/utils"
	"strings"
)

type AuthController struct {
}

func (auth *AuthController) Login(c *gin.Context) {
	var loginForm dtos.LoginForm

	if validationErr := c.ShouldBindJSON(&loginForm); validationErr != nil {
		utils.LogError("BindJSON AuthController", validationErr)
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
	path := c.Request.Header["Path"]
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
	c.Header("x-user-role", jwtClaims.Role)

	if CheckRole(path[0], jwtClaims.Role) {
		c.JSON(http.StatusOK, dtos.Response{Message: "success", Data: nil})
		return
	}

	c.JSON(http.StatusForbidden, gin.H{
		"error": "Forbidden",
	})
}

func CheckRole(path string, role string) bool {
	switch role {
	case "admin":
		for _, value := range middlewares.GetSecurityRouters().Admin {
			if strings.Contains(path, value) {
				return true
			}
		}
		break
	case "user":
		for _, value := range middlewares.GetSecurityRouters().User {
			if strings.Contains(path, value) {
				return true
			}
		}
		break
	}
	return false
}

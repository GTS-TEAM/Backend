package main

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"next/controllers"
	"next/models"
	"os"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error: failed to load the env file")
	}

	if os.Getenv("ENV") == "PRODUCTION" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	models.Init()

	api := r.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Hello World",
			})
		})
		authGroup := api.Group("/auth")
		{
			auth := new(controllers.AuthController)
			authGroup.POST("/login", auth.Login)
			authGroup.POST("/register", auth.Register)
		}
	}

	r.Run(":8080") // listen and serve on
}

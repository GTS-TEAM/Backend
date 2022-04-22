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

func TokenAuthMiddleware() gin.HandlerFunc {
	t := &models.Token{}
	return func(c *gin.Context) {
		t.TokenValid(c)
		c.Next()
	}
}

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
			authGroup.POST("/refresh-token", auth.RefreshToken)
		}
		userGroup := api.Group("/user")
		{
			user := new(controllers.UserController)
			userGroup.GET("/", TokenAuthMiddleware(), user.Get)
			//userGroup.GET("/", user.GetAll)
			//userGroup.GET("/:id", user.Get)
			//userGroup.POST("/", user.Create)
			//userGroup.PUT("/:id", user.Update)
			//userGroup.DELETE("/:id", user.Delete)
		}
	}

	r.Run(":8080") // listen and serve on
}

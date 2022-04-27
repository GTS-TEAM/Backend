package main

import (
	"fmt"
	"github.com/gin-contrib/gzip"
	static "github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/hpcloud/tail"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"io"
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

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uu := uuid.NewV4()
		c.Writer.Header().Set("X-Request-Id", uu.String())
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
	os.Mkdir("logs", 0777)
	logFile, _ := os.Create("logs/server.log")
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)

	server := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		go func() {
			t, err := tail.TailFile("logs/server.log", tail.Config{Follow: true})
			if err != nil {
				log.Fatal(err)
			}
			for line := range t.Lines {
				server.BroadcastToNamespace("/", "some", line.Text)
			}
		}()
		return nil
	})
	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		fmt.Println("Somebody just close the connection ")
	})

	r := gin.Default()

	r.Use(static.Serve("/", static.LocalFile("./public", true)))
	r.Use(CORSMiddleware())
	r.Use(RequestIDMiddleware())
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
			//userGroup.GET("/", user.GetProductsByCategory)
			//userGroup.GET("/:id", user.Gets)
			//userGroup.POST("/", user.Create)
			//userGroup.PUT("/:id", user.Update)
			//userGroup.DELETE("/:id", user.Delete)
		}
		productGroup := api.Group("/product")
		{
			product := new(controllers.ProductController)
			productGroup.GET("/", TokenAuthMiddleware(), product.GetProductsByCategory)
			productGroup.GET("/:id", TokenAuthMiddleware(), product.GetProductsByCategory)
			//productGroup.GET("/:id", product.Gets)
			productGroup.POST("/", TokenAuthMiddleware(), product.Create)
			//productGroup.PUT("/:id", product.Update)
			//productGroup.DELETE("/:id", product.Delete)
		}
		categoryGroup := api.Group("/category")
		{
			category := new(controllers.CategoryController)
			categoryGroup.GET("/", category.GetAll)
			//categoryGroup.GET("/:id", category.Gets)
			categoryGroup.POST("/", category.Create)
			//categoryGroup.PUT("/:id", category.Update)
			//categoryGroup.DELETE("/:id", category.Delete)
		}
	}
	go server.Serve()
	defer server.Close()
	r.GET("/socket.io/", gin.WrapH(server))
	r.POST("/socket.io/*any", gin.WrapH(server))
	// Method 2 using server.ServerHTTP(Writer, Request) and also you can simply this by using gin.WrapH

	r.Run(":8080")
}

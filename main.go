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
	r.RedirectTrailingSlash = true
	r.Use(static.Serve("/", static.LocalFile("./public", true)))
	r.Use(CORSMiddleware())
	r.Use(RequestIDMiddleware())
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	models.Init()

	api := r.Group("/api")
	{
		auth := new(controllers.AuthController)
		api.POST("/auth/login", auth.Login)
		api.POST("/auth/register", auth.Register)
		api.POST("/auth/refresh-token", auth.RefreshToken)

		user := new(controllers.UserController)
		api.GET("/user", TokenAuthMiddleware(), user.Get)
		//userGroup.GET("/", user.GetProductsByCategory)
		//userGroup.GET("/:id", user.Gets)
		//userGroup.POST("/", user.Create)
		//userGroup.PUT("/:id", user.Update)
		//userGroup.DELETE("/:id", user.Delete)
		product := new(controllers.ProductController)
		api.GET("/product", product.GetProductsByCategory)
		api.GET("/product/:id", product.GetProductById)
		api.POST("/product", TokenAuthMiddleware(), product.Create)
		api.PUT("/product/:id", TokenAuthMiddleware(), product.Update)
		api.DELETE("/product/:id", TokenAuthMiddleware(), product.Delete)

		api.GET("/product/reviews/:id", product.GetReviews)
		api.POST("/product/reviews", TokenAuthMiddleware(), product.CreateReviews)

		category := new(controllers.CategoryController)
		api.GET("/category", category.GetAll)
		api.POST("/category", TokenAuthMiddleware(), category.Create)
		//api.POST("/category", category.Create)
		//categoryGroup.PUT("/:id", category.Update)
		//categoryGroup.DELETE("/:id", category.Delete)

		metadata := new(controllers.MetadataController)
		api.GET("/metadata", metadata.GetAll)
		api.POST("/metadata", metadata.Create)
		api.PUT("/metadata/:id", metadata.Update)

		variant := new(controllers.VariantController)
		api.GET("/variant", variant.Get)
		//api.POST("/variant", variant.Create)

		stock := new(controllers.StockController)
		//api.GET("/stock", TokenAuthMiddleware(), stock.Get)
		api.PATCH("/stock", TokenAuthMiddleware(), stock.Update)
		api.GET("/stock", stock.Get)
	}

	go server.Serve()
	defer server.Close()

	r.GET("/socket.io/", gin.WrapH(server))
	r.POST("/socket.io/*any", gin.WrapH(server))

	r.Run(":8080")
}

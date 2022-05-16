package main

import (
	"github.com/gin-contrib/gzip"
	static "github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"io"
	"log"
	"next/controllers"
	"next/middlewares"
	"next/models"
	"os"
)

func main() {

	if os.Getenv("ENV") == "PRODUCTION" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("error: failed to load the env file")
		}
	}
	os.Mkdir("logs", 0777)
	logFile, _ := os.Create("logs/server.log")
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)

	//server := socketio.NewServer(nil)
	//
	//server.OnConnect("/", func(s socketio.Conn) error {
	//	s.SetContext("")
	//	go func() {
	//		t, err := tail.TailFile("logs/server.log", tail.Config{Follow: true})
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//		for line := range t.Lines {
	//			server.BroadcastToNamespace("/", "some", line.Text)
	//		}
	//	}()
	//	return nil
	//})
	//server.OnDisconnect("/", func(s socketio.Conn, msg string) {
	//	fmt.Println("Somebody just close the connection ")
	//})

	r := gin.Default()
	r.RedirectTrailingSlash = true
	r.Use(static.Serve("/", static.LocalFile("./public", true)))
	r.Use(middlewares.CORSMiddleware())
	r.Use(middlewares.RequestIDMiddleware())
	r.Use(gin.Logger())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	//r.Use(gin.Recovery())

	models.Init()

	api := r.Group("/api")
	{
		auth := new(controllers.AuthController)
		api.GET("/authorize", auth.Authorize)
		api.POST("/auth/login", auth.Login)
		api.POST("/auth/register", auth.Register)
		api.POST("/auth/refresh-token", auth.RefreshToken)

		user := new(controllers.UserController)
		api.GET("/user", user.Get)
		//userGroup.GET("/", user.GetProductsByCategory)
		//userGroup.GET("/:id", user.Gets)
		//userGroup.POST("/", user.Create)
		//userGroup.PUT("/:id", user.Update)
		//userGroup.DELETE("/:id", user.Delete)
		product := new(controllers.ProductController)
		api.GET("/product", product.GetProductsByCategory)
		api.GET("/product/:id", product.GetProductById)
		api.POST("/product", product.Create)
		api.PUT("/product/:id", product.Update)
		api.DELETE("/product/:id", product.Delete)

		api.GET("/product/reviews/:id", product.GetReviews)
		api.POST("/product/reviews", product.CreateReviews)

		category := new(controllers.CategoryController)
		api.GET("/category", category.GetAll)
		api.GET("/category/count-products/:id", category.GetCountProductOfCategory)

		api.POST("/category", category.Create)
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
		//api.GET("/stock", (), stock.Get)
		api.PATCH("/stock", stock.Update)
		api.GET("/stock", stock.Get)
	}

	//go server.Serve()
	//defer server.Close()
	//
	//r.GET("/socket.io/", gin.WrapH(server))
	//r.POST("/socket.io/*any", gin.WrapH(server))

	r.Run(":8080")
}

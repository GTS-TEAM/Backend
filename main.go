package main
import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
 "github.com/gin-gonic/gin"
)
func main() {

	dsn := "host=localhost user=postgres password=postgres dbname=next_db port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	_, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Hello World!")
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on
}
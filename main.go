package main
import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
func main() {


	dsn := "host=localhost user=postgres password=postgres dbname=next_db port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	_, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Hello World!")
}
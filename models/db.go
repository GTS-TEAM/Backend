package models

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"sync"
)

type MyDb struct {
	*gorm.DB
}

var (
	once sync.Once
	db   *MyDb
)

//Init ...
func Init() {

	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	err := ConnectDB(dbinfo)
	if err != nil {
		return
	}
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
		fmt.Printf("error while creating MyDb extension 'uuid-ossp': %s\n", err)
	}
	err = db.AutoMigrate(&User{}, &Token{}, &Product{}, &Category{}, &Review{})
	if err != nil {
		log.Println(err)
	}
}

//ConnectDB ...,
func ConnectDB(dataSourceName string) error {
	var err error
	var gdb *gorm.DB
	once.Do(func() {
		gdb, err = gorm.Open(postgres.New(postgres.Config{
			DSN:                  dataSourceName,
			PreferSimpleProtocol: true,
		}), &gorm.Config{})
		if err != nil {
			log.Println("error while connecting to database: ", err)
			return
		}
	})
	db = &MyDb{gdb}
	return nil
}

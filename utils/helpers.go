package utils

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"next/dtos"
	"strconv"
)

func LogError(in string, err error) {
	log.Println("\n[ERROR] " + "[" + in + "] " + err.Error())
}

func LogInfo(in string, message string) {
	log.Println("\n[INFO] " + "[" + in + "] " + message)
}

func BindStruct(obj interface{}, to interface{}) error {
	marshal, err := json.Marshal(obj)

	if err != nil {
		return err
	}
	return json.Unmarshal(marshal, to)
}

func GetProductFilter(c *gin.Context) (filter dtos.ProductFilter) {
	filter.MaxPrice, _ = strconv.ParseFloat(c.Query("max_price"), 32)
	filter.MinPrice, _ = strconv.ParseFloat(c.Query("min_price"), 32)
	filter.MinRating, _ = strconv.ParseFloat(c.Query("min_rating"), 32)
	filter.Name = c.Query("name")
	return
}

func GetCustomersFilter(c *gin.Context) (filter dtos.CustomerFilter) {
	filter.Name = c.Query("name")
	filter.Email = c.Query("email")
	filter.Phone = c.Query("phone")
	filter.Status = c.Query("status")
	return
}

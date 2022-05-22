package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"next/dtos"
	"next/models"
	"next/utils"
)

/* Utils functions */
func GetQueryChangeStatus(c *gin.Context) error {
	var dto dtos.ChangeStatusRequest
	dto.CustomerID = c.Query("customer_id")
	dto.Status = c.Query("status")
	if dto.CustomerID == "" || dto.Status == "" {
		return errors.New("Invalid params")
	}
	if !checkValidStatus(dto.Status) {
		return errors.New("Invalid status")
	}

	return nil
}

func checkValidStatus(status string) bool {
	if status == "active" || status == "block" {
		return true
	}
	return false
}

func getUserID(c *gin.Context) (userID string) {
	//MustGet returns the value for the given key if it exists, otherwise it panics.
	// get x-user-id from header
	return c.Request.Header.Get("x-user-id")
}

/* Controller functions */
type UserController struct {
}

func (u *UserController) GetCustomers(c *gin.Context) {
	user := models.User{}

	paging := models.GeneratePaginationFromRequest(c)
	filter := utils.GetCustomersFilter(c)

	customers, err := user.GetCustomers(filter, paging)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, dtos.Response{
		Data: customers,
	})
}

func (u *UserController) ChangeCustomerStatus(c *gin.Context) {
	err := GetQueryChangeStatus(c)
	request := dtos.ChangeStatusRequest{}
	user := models.User{}
	if err != nil {
		c.JSON(406, gin.H{"error": err.Error()})
		return
	}
	err = user.ChangeStatus(request)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, dtos.Response{
		Message: "success",
		Data:    nil,
	})
}

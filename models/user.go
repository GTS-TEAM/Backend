package models

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"next/dtos"
)

type User struct {
	BaseModel
	Name     string    `json:"name"`
	Email    *string   `gorm:"unique" json:"email"`
	Password string    `json:"-"`
	Role     string    `json:"role"`
	Products []Product `json:"products,omitempty"`
	Phone    string    `json:"phone" gorm:"unique;null"`
	Status   string    `json:"status" gorm:"default:'active'"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	h, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(h)
	return nil
}

func (u *User) Login(dto dtos.LoginForm) (res map[string]interface{}, err error) {
	user := User{}
	if err := db.First(&user, "email = ? and role = ?", dto.Email, dto.Role).Error; err == gorm.ErrRecordNotFound {
		return nil, errors.New("Email not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password)); err != nil {
		return nil, errors.New("Password not match")
	}

	token := Token{}

	authToken, err := token.GenerateAuthToken(user.ID)

	if err != nil {
		fmt.Printf("Error generating token: %v", err)
		return nil, errors.New("Failed to generate auth token")
	}

	return map[string]interface{}{
		"user":  user,
		"token": authToken,
	}, nil
}

func (u *User) Register(dto dtos.RegisterForm) (err error) {
	user := User{}
	if err := db.First(&user, "email = ?", dto.Email).Error; err == nil {
		return errors.New("Email already exists")
	}

	user = User{
		Name:     dto.Name,
		Email:    &dto.Email,
		Password: dto.Password,
		Role:     dto.Role,
	}

	if err := db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (u *User) GetUserByID(id uint) (user User, err error) {
	if err := db.First(&user, id).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (u *User) GetName(userId string) (user User, err error) {
	if err = db.First(&user, "id = ?", userId).Error; err != nil {
		return user, err
	}
	fmt.Println("user", user)
	return user, nil
}

func (u *User) GetUserByEmail(email string) (user User, err error) {
	if err := db.First(&user, "email = ?", email).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (u *User) GetCustomers(filter dtos.CustomerFilter, pagination Pagination) (users []User, err error) {
	qb := db.Where("role = ?", "customer")

	if filter.Name != "" {
		qb = qb.Where("name ILIKE ?", "%"+filter.Name+"%")
	}
	if filter.Email != "" {
		qb = qb.Where("email ILIKE ?", "%"+filter.Email+"%")
	}
	if filter.Phone != "" {
		qb = qb.Where("phone LIKE ?", "%"+filter.Phone+"%")
	}
	if filter.Status != "" {
		qb = qb.Where("status = ?", filter.Status)
	}

	if err = qb.Offset(pagination.Page).Limit(pagination.Limit).Find(&users).Error; err != nil {
		panic(err)
		return users, err
	}

	return users, nil
}

func (u *User) ChangeStatus(request dtos.ChangeStatusRequest) (err error) {
	user := User{}
	if err = db.Debug().First(&user, "id = ?", request.CustomerID).Error; err != nil {
		return err
	}

	user.Status = request.Status

	if err = db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

//func (u *User) Delete(id string) (err error) {
//	user := User{}
//	if err = db.First(&user, id).Error; err != nil {
//		return err
//	}
//
//	if err = db.Delete(&user).Error; err != nil {
//		return err
//	}
//
//	return nil
//}

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
	Name     string  `json:"name"`
	Email    *string `gorm:"unique" json:"email"`
	Password string  `json:"-"`
	Role     string  `json:"role"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	h, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(h)
	return nil
}

func (u *User) Login(dto dtos.LoginForm) (res map[string]interface{}, err error) {
	user := User{}
	if err := db.First(&user, "email = ?", dto.Email).Error; err == gorm.ErrRecordNotFound {
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
		Role:     "user",
	}

	if err := db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

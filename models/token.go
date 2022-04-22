package models

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Token struct {
	BaseModel
	Token  string    `json:"token"`
	Type   string    `json:"type"`
	UserID uuid.UUID `json:"user"`
}

type AuthToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type JwtCustomClaim struct {
	UserID    uuid.UUID `json:"user_id"`
	IsValid   bool      `json:"is_valid"`
	Role      string    `json:"role"`
	TokenType string    `json:"token_type"`
	jwt.StandardClaims
}

const (
	AccessTokenType  = "access"
	RefreshTokenType = "refresh"
)

func (t *Token) GenerateToken(UserID uuid.UUID, tokenType string) (string, error) {

	var user User
	if err := db.First(&user, "id = ?", UserID.String()).Error; err != nil {
		return "", err
	}

	claims := &JwtCustomClaim{
		UserID,
		true,
		user.Role,
		tokenType,
		jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(1, 0, 0).Unix(),
			Issuer:    "go-jwt",
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tk, err := token.SignedString(getSecretKey())
	if err != nil {
		panic(err)
	}
	return tk, nil
}

func (t *Token) GenerateAuthToken(UserID uuid.UUID) (*AuthToken, error) {
	accessToken, err := t.GenerateToken(UserID, AccessTokenType)
	refreshToken, err := t.GenerateToken(UserID, RefreshTokenType)

	db.Create(&Token{
		Token:  refreshToken,
		Type:   RefreshTokenType,
		UserID: UserID,
	})

	if err != nil {
		return nil, err
	}

	return &AuthToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (t *Token) RefreshToken(token string) (*AuthToken, error) {
	fmt.Println("RefreshToken: ", token)
	userId, err := t.VerifyToken(token)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fromString, err := uuid.FromString(userId)
	if err != nil {
		return nil, err
	}

	return t.GenerateAuthToken(fromString)
}

func (t *Token) VerifyToken(token string) (string, error) {
	jwtPayload, err := jwt.Parse(token, func(t_ *jwt.Token) (interface{}, error) {
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method %v ", t_.Header["alg"])
		}
		return getSecretKey(), nil
	})

	if err != nil {
		return "", err
	}

	userId := jwtPayload.Claims.(jwt.MapClaims)["user_id"].(string)
	fmt.Println("User Id: ", userId)

	if err = db.First(&Token{}, "token = ? AND user_id = ?", token, userId).Error; err != nil {
		fmt.Println("VerifyToken Err: ", err)
		return "", errors.New("Invalid Token")
	}

	return userId, nil
}

func getSecretKey() []byte {
	return []byte("secret")
}

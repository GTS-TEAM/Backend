package models

import (
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
	UserID  uuid.UUID `json:"user_id"`
	IsValid bool      `json:"is_valid"`
	Role    string    `json:"role"`
	jwt.StandardClaims
}

const (
	AccessTokenType  = "access"
	RefreshTokenType = "refresh"
)

func (t *Token) GenerateToken(UserID uuid.UUID) (string, error) {

	var user User
	if err := db.First(&user, "id = ?", UserID.String()).Error; err != nil {
		return "", err
	}

	claims := &JwtCustomClaim{
		UserID,
		true,
		user.Role,
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
	accessToken, err := t.GenerateToken(UserID)
	refreshToken, err := t.GenerateToken(UserID)

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

//func (t *Token) ValidateToken(token string) (*jwt.Token, error) {
//	return jwt.Parse(token, func(t_ *jwt.Token) (interface{}, error) {
//		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
//			return nil, fmt.Errorf("Unexpected signing method %v ", t_.Header["alg"])
//		}
//		return getSecretKey(), nil
//	})
//}

func getSecretKey() []byte {
	return []byte("secret")
}

package models

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strings"
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

func (t *Token) GenerateToken(UserID uuid.UUID, tokenType string, exp int64) (string, error) {
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
			ExpiresAt: exp,
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

func (t *Token) TokenValid(c *gin.Context) {

	token := t.ExtractToken(c.Request)

	userId, err := t.VerifyToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err != nil {
		//Token does not exists in Redis (User logged out or expired)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Please login first"})
		return
	}

	//To be called from GetUserID()
	c.Set("userID", userId)
}

func (t *Token) ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func (t *Token) GenerateAuthToken(UserID uuid.UUID) (*AuthToken, error) {
	accessToken, err := t.GenerateToken(UserID, AccessTokenType, time.Now().Add(time.Minute*30).Unix())
	refreshToken, err := t.GenerateToken(UserID, RefreshTokenType, time.Now().Add(time.Hour*24*30).Unix())

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

func (t *Token) VerifyToken(token string) (uuid.UUID, error) {
	jwtPayload, err := jwt.Parse(token, func(t_ *jwt.Token) (interface{}, error) {
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method %v ", t_.Header["alg"])
		}
		return getSecretKey(), nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	userId := jwtPayload.Claims.(jwt.MapClaims)["user_id"].(string)
	return uuid.FromString(userId)
}

func (t *Token) ValidateTokenRefreshToken(refreshToken string) (*AuthToken, error) {
	userId, err := t.VerifyToken(refreshToken)
	if err != nil {
		return nil, err
	}

	if err = db.First(&Token{}, "token = ? AND user_id = ?", refreshToken, userId).Error; err != nil {
		return nil, errors.New("Invalid Token")
	}

	return t.GenerateAuthToken(userId)
}

func getSecretKey() []byte {
	return []byte("secret")
}

package models

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"next/utils"
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

type AuthorizationResponse struct {
	XUserID uuid.UUID `json:"x-user-id"`
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

	jwtClaims, err := t.VerifyToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err != nil {
		//Token does not exists in Redis (Customer logged out or expired)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Token is valid",
		"claims":  jwtClaims,
	})
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
	accessToken, err := t.GenerateToken(UserID, AccessTokenType, time.Now().Add(time.Hour*24*30).Unix())
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

func (t *Token) VerifyToken(token string) (jwtClaims JwtCustomClaim, err error) {
	jwtPayload, err := jwt.Parse(token, func(t_ *jwt.Token) (interface{}, error) {
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method %v ", t_.Header["alg"])
		}
		return getSecretKey(), nil
	})

	if err != nil {
		return jwtClaims, err
	}

	err = utils.BindStruct(jwtPayload.Claims.(jwt.MapClaims), &jwtClaims)
	if err != nil {
		fmt.Printf("Error Bind struct: %v\n", err)
	}

	if jwtClaims.UserID == uuid.Nil {
		return jwtClaims, errors.New("Invalid Token")
	}

	if !jwtClaims.IsValid {
		return jwtClaims, errors.New("Invalid Token")
	}

	return jwtClaims, nil
}

func (t *Token) ValidateTokenRefreshToken(refreshToken string) (*AuthToken, error) {
	jwtClaims, err := t.VerifyToken(refreshToken)
	if err != nil {
		return nil, err
	}

	if err = db.First(&Token{}, "token = ? AND user_id = ?", refreshToken, jwtClaims.UserID).Error; err != nil {
		return nil, errors.New("Invalid Token")
	}

	return t.GenerateAuthToken(jwtClaims.UserID)
}

func getSecretKey() []byte {
	return []byte("secret")
}

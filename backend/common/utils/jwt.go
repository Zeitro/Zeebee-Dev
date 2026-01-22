package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"zeebee-backend/domains/zeebee/constants"
)

const (
	JwtTokenCtxKey = "jwt_token"
)

// JWTClaims JWT声明结构
type JWTClaims struct {
	JwtUser
	jwt.RegisteredClaims
}

type JwtUser struct {
	UserId             int64              `json:"user_id"`
	OrganizationId     int64              `json:"organization_id"`
	OrganizationUserId int64              `json:"organization_user_id"`
	Username           string             `json:"username"`
	Role               []string           `json:"role"`
	UserType           constants.UserType `json:"user_type"`
}

// GenerateToken 生成JWT token
func GenerateToken(user JwtUser, secretKey string, expireSeconds int64) (string, error) {
	now := time.Now()
	expireTime := now.Add(time.Duration(expireSeconds) * time.Second)

	claims := JWTClaims{
		JwtUser: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "zeebee",
			Subject:   fmt.Sprintf("%v", user.UserId),
			ID:        uuid.New().String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// ParseToken 解析JWT token
func ParseToken(tokenString, secretKey string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

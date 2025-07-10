package middleware

import (
	"context"
	"net/http"
	"strings"

	"zeebee-backend/common/response"
	"zeebee-backend/common/utils"
)

type AuthMiddleware struct {
	secretKey string
}

func NewAuthMiddleware(secretKey string) *AuthMiddleware {
	return &AuthMiddleware{
		secretKey: secretKey,
	}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 跳过健康检查等公开接口
		if r.URL.Path == "/health" || r.URL.Path == "/metrics" {
			next(w, r)
			return
		}

		// 获取Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Unauthorized(w, "Authorization header is required")
			return
		}

		// 检查Bearer token格式
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(w, "Invalid authorization header format")
			return
		}

		token := parts[1]

		// 解析JWT token
		claims, err := utils.ParseToken(token, m.secretKey)
		if err != nil {
			response.Unauthorized(w, "Invalid token")
			return
		}

		// 将用户信息存储到请求上下文中
		ctx := r.Context()
		ctx = context.WithValue(ctx, "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "username", claims.Username)
		ctx = context.WithValue(ctx, "role", claims.Role)

		next(w, r.WithContext(ctx))
	}
}

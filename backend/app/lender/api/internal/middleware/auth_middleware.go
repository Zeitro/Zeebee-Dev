package middleware

import (
	"context"
	"net/http"
	"strings"

	"zeebee-backend/common/constants"
	"zeebee-backend/common/response"
	"zeebee-backend/common/utils"
)

type AuthMiddleware struct {
	secretKey string
}

func NewAuthMiddleware(secretKey string) *AuthMiddleware {
	return &AuthMiddleware{secretKey: secretKey}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get user info from header
		token := r.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")
		if token == "" {
			response.Unauthorized(w, "Unauthorized")
			return
		}
		claims, err := utils.ParseToken(token, m.secretKey)
		if err != nil {
			response.Unauthorized(w, "Invalid token")
			return
		}
		if claims.UserType != constants.UserTypeLender {
			response.Unauthorized(w, "Invalid user type")
			return
		}

		reqCtx := r.Context()
		nextCtx := context.WithValue(reqCtx, constants.UserClaimCtxKey, claims)
		nextCtx = context.WithValue(nextCtx, constants.UserIdCtxKey, claims.UserId)
		// nextCtx = context.WithValue(nextCtx, constants.LenderUserCompanyIdKey, claims.OrganizationId)
		newReq := r.WithContext(nextCtx)

		next(w, newReq)
	}
}

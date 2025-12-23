package middleware

import (
	"context"
	"net/http"
	"strings"

	"google.golang.org/grpc/metadata"

	"zeebee-backend/common/response"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
)

type HttpAuthMiddleware struct {
	secretKey string
	userType  constants.UserType
}

func NewHttpAuthMiddleware(secretKey string, userType constants.UserType) *HttpAuthMiddleware {
	return &HttpAuthMiddleware{secretKey: secretKey, userType: userType}
}

func (m *HttpAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
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
			response.Unauthorized(w, "Invalid Token")
			return
		}
		if claims.UserType != m.userType {
			response.Unauthorized(w, "Invalid User Type")
			return
		}

		reqCtx := r.Context()
		nextCtx := context.WithValue(reqCtx, utils.UserClaimCtxKey, claims)
		nextCtx = context.WithValue(nextCtx, utils.UserIdCtxKey, claims.UserId)
		nextCtx = metadata.AppendToOutgoingContext(nextCtx, utils.JwtTokenCtxKey, token)
		newReq := r.WithContext(nextCtx)

		next(w, newReq)
	}
}

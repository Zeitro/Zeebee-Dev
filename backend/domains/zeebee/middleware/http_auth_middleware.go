package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/pathvar"
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

// getRoutePattern converts a concrete path to a parameterized route pattern
// e.g., /api/v1/resource/7 -> /api/v1/resource/:id
// Uses go-zero's pathvar to get actual route parameters and reconstruct the pattern
func getRoutePattern(r *http.Request) string {
	path := r.URL.Path
	vars := pathvar.Vars(r)

	// If no path variables, return the path as-is
	if len(vars) == 0 {
		return path
	}

	// Reconstruct the route pattern by replacing parameter values with parameter names
	pattern := path
	for key, value := range vars {
		// Replace the value in the path with :key pattern
		// e.g., /api/v1/resource/7 -> /api/v1/resource/:id
		pattern = strings.Replace(pattern, "/"+value, "/:"+key, 1)
	}

	return pattern
}

func (m *HttpAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		routePattern := getRoutePattern(r)
		method := r.Method
		logx.Infof("route: %s, method: %s", routePattern, method)

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
		if claims.UserType != m.userType {
			response.Unauthorized(w, "Invalid user type")
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

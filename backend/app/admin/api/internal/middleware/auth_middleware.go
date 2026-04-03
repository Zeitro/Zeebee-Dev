package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"

	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/response"
	"zeebee-backend/common/routex"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
)

type AuthMiddleware struct {
	userType constants.UserType

	authRpcClient admin.AuthClient
}

func NewAuthMiddleware(userType constants.UserType, authRpcClient admin.AuthClient) *AuthMiddleware {
	return &AuthMiddleware{userType: userType, authRpcClient: authRpcClient}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCtx := r.Context()
		routePattern := routex.GetGoZeroRoutePattern(r)
		method := r.Method
		// logx.Infof("route: %s, method: %s", routePattern, method)

		// get user info from header
		token := r.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")
		nextCtx := context.WithValue(reqCtx, utils.JwtTokenCtxKey, token)
		nextCtx = metadata.AppendToOutgoingContext(nextCtx, utils.JwtTokenCtxKey, token)

		rpcResp, err := m.authRpcClient.AuthCheck(nextCtx, &admin.AuthCheckRequest{
			Data: &admin.AuthCheckRequestData{
				Method: method,
				Path:   routePattern,
				Type:   string(m.userType),
			},
		})
		if err != nil {
			logx.WithContext(reqCtx).
				WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).
				Errorf("Auth check failed: %v", err)
			response.BadRequest(w, "Authentication failed")
			return
		}

		if rpcResp.Data == nil {
			response.BadRequest(w, "Invalid authentication response")
			return
		}
		result := rpcResp.Data
		if result.AuthenticationRequired && !result.Authenticated {
			response.Unauthorized(w, "Not Login")
			return
		}

		if result.AuthorizationRequired && !result.Authorized {
			response.Forbidden(w, "No permission")
			return
		}

		nextCtx = metadata.AppendToOutgoingContext(nextCtx, constants.DataScopeCtxKey, result.DataScope)

		newReq := r.WithContext(nextCtx)

		next(w, newReq)
	}
}

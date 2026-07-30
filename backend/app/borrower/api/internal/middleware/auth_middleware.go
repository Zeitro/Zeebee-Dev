package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"

	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/response"
	"zeebee-backend/common/routex"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
)

type AuthMiddleware struct {
	userType constants.UserType

	authRpcClient borrower.AuthClient
}

func NewAuthMiddleware(userType constants.UserType, authRpcClient borrower.AuthClient) *AuthMiddleware {
	return &AuthMiddleware{
		userType:      userType,
		authRpcClient: authRpcClient,
	}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCtx := r.Context()
		routePattern := routex.GetGoZeroRoutePattern(r)
		method := r.Method

		// get user info from header
		token := r.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")
		nextCtx := context.WithValue(reqCtx, utils.JwtTokenCtxKey, token)
		nextCtx = metadata.AppendToOutgoingContext(nextCtx, utils.JwtTokenCtxKey, token)

		rpcResp, err := m.authRpcClient.AuthCheck(nextCtx, &proto.AuthCheckRequest{
			Data: &proto.AuthCheckRequestData{
				Method: method,
				Path:   routePattern,
				Type:   string(m.userType),
			},
		})
		if err != nil {
			logx.WithContext(reqCtx).
				WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).
				Errorf("Auth check failed: %v", err)
			response.BadRequest(w, "authentication failed")
			return
		}

		if rpcResp.Data == nil {
			response.BadRequest(w, "invalid auth response")
			return
		}
		result := rpcResp.Data
		if result.AuthenticationRequired && !result.Authenticated {
			response.Unauthorized(w, "not login")
			return
		}

		if result.AuthorizationRequired && !result.Authorized {
			response.Forbidden(w, "no permission")
			return
		}

		nextCtx = metadata.AppendToOutgoingContext(nextCtx, constants.DataScopeCtxKey, result.DataScope)

		newReq := r.WithContext(nextCtx)

		next(w, newReq)
	}
}

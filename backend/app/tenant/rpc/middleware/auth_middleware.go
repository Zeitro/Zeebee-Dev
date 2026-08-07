package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/response"
	"zeebee-backend/common/routex"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
)

// AuthChecker is the authz surface the middleware needs. tenant.AuthzCheck
// satisfies it; a local per-backend Auth RPC client does too (same proto shape).
type AuthChecker interface {
	AuthCheck(ctx context.Context, in *proto.AuthCheckRequest, opts ...grpc.CallOption) (*proto.AuthCheckResponse, error)
}

// AuthMiddleware enforces identity + authorization for protected routes. It
// forwards the caller's JWT to the auth checker (tenant.AuthzCheck) and, before
// any decision, verifies locally that the claim's userType matches this
// service's platform type. Checker failures are mapped to business codes in the
// standard response body (HTTP stays 200, matching the frontend convention);
// callers tell a network outage (503) from an auth denial (401/403) by code.
type AuthMiddleware struct {
	userType  constants.UserType
	checker   AuthChecker
	jwtSecret string
}

func NewAuthMiddleware(userType constants.UserType, checker AuthChecker, jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{userType: userType, checker: checker, jwtSecret: jwtSecret}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCtx := r.Context()
		routePattern := routex.GetGoZeroRoutePattern(r)
		method := r.Method

		// get user info from header
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if token != "" {
			claim, err := utils.ParseToken(token, m.jwtSecret)
			if err != nil {
				// response.Unauthorized(w, "invalid token")
				response.BadRequest(w, "authentication failed")
				return
			}
			if claim.UserType != m.userType {
				response.Unauthorized(w, "user type mismatch")
				return
			}
		}
		nextCtx := context.WithValue(reqCtx, utils.JwtTokenCtxKey, token)
		nextCtx = metadata.AppendToOutgoingContext(nextCtx,
			utils.JwtTokenCtxKey, token,
			utils.UserTypeCtxKey, string(m.userType),
		)

		rpcResp, err := m.checker.AuthCheck(nextCtx, &proto.AuthCheckRequest{
			Data: &proto.AuthCheckRequestData{
				Method: method,
				Path:   routePattern,
				Type:   string(m.userType),
			},
		})
		if err != nil {
			switch status.Code(err) {
			case codes.Unavailable, codes.DeadlineExceeded:
				// tenant unreachable: judge availability from the HTTP status.
				// Degradation (decision cache) hooks in here per NFR-002.
				logx.WithContext(reqCtx).
					WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).
					Errorf("auth service unavailable: %v", err)
				response.Error(w, http.StatusServiceUnavailable, "auth service unavailable")
			case codes.Unauthenticated:
				response.Unauthorized(w, "not login")
			case codes.PermissionDenied:
				response.Forbidden(w, "no permission")
			default:
				logx.WithContext(reqCtx).
					WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).
					Errorf("Auth check failed: %v", err)
				response.Error(w, http.StatusInternalServerError, "authentication failed")
			}
			return
		}

		if rpcResp.Data == nil {
			response.Error(w, http.StatusInternalServerError, "invalid authentication response")
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
		next(w, r.WithContext(nextCtx))
	}
}

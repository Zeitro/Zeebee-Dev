package middleware

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	tenant "zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/utils"
	authdomain "zeebee-backend/domains/zeebee/auth"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type bodyCode struct {
	Code int `json:"code"`
}

type fakeChecker struct {
	called bool
	ctx    context.Context
	in     *proto.AuthCheckRequest
	resp   *proto.AuthCheckResponse
	err    error
}

type authCheckServer struct {
	tenant.UnimplementedAuthServer
	auth *authdomain.Auth
}

func (s authCheckServer) AuthCheck(ctx context.Context, _ *proto.AuthCheckRequest) (*proto.AuthCheckResponse, error) {
	_, err := s.auth.GetUserClaimFromCtx(ctx, constants.UserTypeLender)
	return nil, err
}

func (f *fakeChecker) AuthCheck(ctx context.Context, in *proto.AuthCheckRequest, opts ...grpc.CallOption) (*proto.AuthCheckResponse, error) {
	f.called = true
	f.ctx = ctx
	f.in = in
	return f.resp, f.err
}

func grantedResp() *proto.AuthCheckResponse {
	return &proto.AuthCheckResponse{Data: &proto.AuthCheckResponseData{
		AuthenticationRequired: true, Authenticated: true,
		AuthorizationRequired: false, Authorized: false,
	}}
}

func runRequest(t *testing.T, mw *AuthMiddleware, token string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/x", nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	mw.Handle(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })(rec, req)
	return rec
}

// bodyCodeOf decodes the standard response body and returns its business code.
// The middleware keeps the frontend convention: HTTP 200 + code in the body.
func bodyCodeOf(t *testing.T, rec *httptest.ResponseRecorder) int {
	t.Helper()
	var b bodyCode
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &b))
	return b.Code
}

// Token verification and user-type matching are delegated to tenant.
func TestAuthMiddleware_MatchingTypePasses(t *testing.T) {
	checker := &fakeChecker{resp: grantedResp()}
	mw := NewAuthMiddleware(constants.UserTypeLender, checker)
	rec := runRequest(t, mw, "tenant-issued-token")

	require.Equal(t, http.StatusOK, rec.Code)
	require.True(t, checker.called)
	md, ok := metadata.FromOutgoingContext(checker.ctx)
	require.True(t, ok)
	require.Equal(t, []string{"tenant-issued-token"}, md.Get(utils.JwtTokenCtxKey))
	require.Equal(t, []string{string(constants.UserTypeLender)}, md.Get(utils.UserTypeCtxKey))
}

func TestAuthMiddleware_InvalidTokenRejectedByTenant(t *testing.T) {
	checker := &fakeChecker{err: status.Error(codes.Unauthenticated, "invalid token")}
	mw := NewAuthMiddleware(constants.UserTypeLender, checker)
	rec := runRequest(t, mw, "invalid-token")

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, http.StatusUnauthorized, bodyCodeOf(t, rec))
	require.True(t, checker.called)
}

func TestAuthMiddleware_GrpcUnauthenticated(t *testing.T) {
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	tenant.RegisterAuthServer(server, authCheckServer{
		auth: authdomain.NewAuth(nil, nil, dto.JwtConfig{AccessSecret: "secret"}, nil, authdomain.AuthTopics{}),
	})
	go func() { _ = server.Serve(listener) }()
	t.Cleanup(func() {
		server.Stop()
		_ = listener.Close()
	})

	conn, err := grpc.DialContext(context.Background(), "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return listener.Dial() }),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	mw := NewAuthMiddleware(constants.UserTypeLender, tenant.NewAuthClient(conn))
	for _, token := range []string{"", "invalid-token"} {
		t.Run(token, func(t *testing.T) {
			rec := runRequest(t, mw, token)

			require.Equal(t, http.StatusOK, rec.Code)
			require.Equal(t, http.StatusUnauthorized, bodyCodeOf(t, rec))
		})
	}
}

// A network failure from the authority surfaces as business code 503 (HTTP 200)
// so callers can tell the tenant service is unreachable by body code.
func TestAuthMiddleware_NetworkUnavailable(t *testing.T) {
	checker := &fakeChecker{err: status.Error(codes.Unavailable, "connection refused")}
	mw := NewAuthMiddleware(constants.UserTypeLender, checker)
	rec := runRequest(t, mw, "tenant-issued-token")

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, http.StatusServiceUnavailable, bodyCodeOf(t, rec))
	require.True(t, checker.called)
}

// A permission denial from the authority surfaces as business code 403.
func TestAuthMiddleware_PermissionDenied(t *testing.T) {
	checker := &fakeChecker{err: status.Error(codes.PermissionDenied, "no permission")}
	mw := NewAuthMiddleware(constants.UserTypeLender, checker)
	rec := runRequest(t, mw, "tenant-issued-token")

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, http.StatusForbidden, bodyCodeOf(t, rec))
	require.True(t, checker.called)
}

func TestAuthMiddleware_ProjectManagementDeniedBeforeHandler(t *testing.T) {
	checker := &fakeChecker{resp: &proto.AuthCheckResponse{Data: &proto.AuthCheckResponseData{
		AuthenticationRequired: true,
		Authenticated:          true,
		AuthorizationRequired:  true,
		Authorized:             false,
	}}}
	mw := NewAuthMiddleware(constants.UserTypeBorrower, checker)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/project_management/workflow/all", nil)
	req.Header.Set("Authorization", "Bearer tenant-issued-token")
	rec := httptest.NewRecorder()
	handled := false

	mw.Handle(func(http.ResponseWriter, *http.Request) { handled = true })(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, http.StatusForbidden, bodyCodeOf(t, rec))
	require.False(t, handled)
	require.True(t, checker.called)
	require.Equal(t, http.MethodGet, checker.in.GetData().GetMethod())
	require.Equal(t, "/api/v1/project_management/workflow/all", checker.in.GetData().GetPath())
	require.Equal(t, string(constants.UserTypeBorrower), checker.in.GetData().GetType())
}

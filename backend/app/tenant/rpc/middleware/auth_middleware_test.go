package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
)

type bodyCode struct {
	Code int `json:"code"`
}

type fakeChecker struct {
	called bool
	resp   *proto.AuthCheckResponse
	err    error
}

func (f *fakeChecker) AuthCheck(ctx context.Context, in *proto.AuthCheckRequest, opts ...grpc.CallOption) (*proto.AuthCheckResponse, error) {
	f.called = true
	return f.resp, f.err
}

func signedToken(t *testing.T, secret string, userType constants.UserType) string {
	t.Helper()
	tok, err := utils.GenerateToken(utils.JwtUser{UserId: 1, OrganizationId: 1, UserType: userType}, secret, 3600)
	require.NoError(t, err)
	return tok
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

// A token whose claim userType differs from the service type is rejected before
// the authority is consulted.
func TestAuthMiddleware_UserTypeMismatchRejected(t *testing.T) {
	checker := &fakeChecker{resp: grantedResp()}
	mw := NewAuthMiddleware(constants.UserTypeLender, checker, "secret")
	rec := runRequest(t, mw, signedToken(t, "secret", constants.UserTypeBorrower))

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, http.StatusUnauthorized, bodyCodeOf(t, rec))
	require.False(t, checker.called)
}

// A token whose claim userType matches the service type passes through.
func TestAuthMiddleware_MatchingTypePasses(t *testing.T) {
	checker := &fakeChecker{resp: grantedResp()}
	mw := NewAuthMiddleware(constants.UserTypeLender, checker, "secret")
	rec := runRequest(t, mw, signedToken(t, "secret", constants.UserTypeLender))

	require.Equal(t, http.StatusOK, rec.Code)
	require.True(t, checker.called)
}

// A token signed with the wrong secret is rejected locally.
func TestAuthMiddleware_InvalidSignatureRejected(t *testing.T) {
	checker := &fakeChecker{resp: grantedResp()}
	mw := NewAuthMiddleware(constants.UserTypeLender, checker, "secret")
	rec := runRequest(t, mw, signedToken(t, "other-secret", constants.UserTypeLender))

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, http.StatusUnauthorized, bodyCodeOf(t, rec))
	require.False(t, checker.called)
}

// A network failure from the authority surfaces as business code 503 (HTTP 200)
// so callers can tell the tenant service is unreachable by body code.
func TestAuthMiddleware_NetworkUnavailable(t *testing.T) {
	checker := &fakeChecker{err: status.Error(codes.Unavailable, "connection refused")}
	mw := NewAuthMiddleware(constants.UserTypeLender, checker, "secret")
	rec := runRequest(t, mw, signedToken(t, "secret", constants.UserTypeLender))

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, http.StatusServiceUnavailable, bodyCodeOf(t, rec))
	require.True(t, checker.called)
}

// A permission denial from the authority surfaces as business code 403.
func TestAuthMiddleware_PermissionDenied(t *testing.T) {
	checker := &fakeChecker{err: status.Error(codes.PermissionDenied, "no permission")}
	mw := NewAuthMiddleware(constants.UserTypeLender, checker, "secret")
	rec := runRequest(t, mw, signedToken(t, "secret", constants.UserTypeLender))

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, http.StatusForbidden, bodyCodeOf(t, rec))
	require.True(t, checker.called)
}

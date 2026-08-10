package logic

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/repo"
)

// newClaimTestUserLogic builds a UserLogic whose jwtConfig signs/verifies with
// the given secret — enough to exercise GetUserClaimFromCtx in isolation.
func newClaimTestUserLogic(t *testing.T, secret string) *UserLogic {
	t.Helper()
	db := setupTestDB(t)
	baseLogic := NewBaseLogic(repo.NewAllRepos(db), nil, nil, dto.JwtConfig{AccessSecret: secret})
	return NewUserLogic(baseLogic, nil)
}

// A token whose embedded userType disagrees with the platform type declared in
// metadata is rejected as invalid — the cross-platform replay guard.
func TestGetUserClaimFromCtx_UserTypeMismatch_InvalidToken(t *testing.T) {
	ul := newClaimTestUserLogic(t, "test-secret")
	tok, err := utils.GenerateToken(utils.JwtUser{
		UserId:   1,
		UserType: constants.UserTypeLender,
	}, "test-secret", 3600)
	require.NoError(t, err)

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		utils.JwtTokenCtxKey, tok,
		utils.UserTypeCtxKey, string(constants.UserTypeBorrower),
	))

	_, err = ul.GetUserClaimFromCtx(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token")
}

// A token whose embedded userType matches the platform type declared in metadata
// is accepted and returns the claim.
func TestGetUserClaimFromCtx_UserTypeMatch_OK(t *testing.T) {
	ul := newClaimTestUserLogic(t, "test-secret")
	tok, err := utils.GenerateToken(utils.JwtUser{
		UserId:   1,
		UserType: constants.UserTypeLender,
	}, "test-secret", 3600)
	require.NoError(t, err)

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		utils.JwtTokenCtxKey, tok,
		utils.UserTypeCtxKey, string(constants.UserTypeLender),
	))

	claim, err := ul.GetUserClaimFromCtx(ctx)
	require.NoError(t, err)
	assert.Equal(t, constants.UserTypeLender, claim.UserType)
}

func TestGetUserClaimFromCtx_InvalidSignature_Unauthenticated(t *testing.T) {
	ul := newClaimTestUserLogic(t, "test-secret")
	tok, err := utils.GenerateToken(utils.JwtUser{UserId: 1, UserType: constants.UserTypeLender}, "other-secret", 3600)
	require.NoError(t, err)

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		utils.JwtTokenCtxKey, tok,
		utils.UserTypeCtxKey, string(constants.UserTypeLender),
	))

	_, err = ul.GetUserClaimFromCtx(ctx)
	require.Error(t, err)
	assert.Equal(t, codes.Unauthenticated, status.Code(err))
}

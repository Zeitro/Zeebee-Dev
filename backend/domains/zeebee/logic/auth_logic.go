package logic

import (
	"context"

	"zeebee-backend/common/rbac"
	"zeebee-backend/domains/zeebee/auth"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
)

// AuthLogic provides authentication business logic
type AuthLogic struct {
	*BaseLogic
	auth *auth.Auth
}

// NewAuthLogic creates a new authentication logic instance
func NewAuthLogic(baseLogic *BaseLogic, jwtConfig dto.JwtConfig, enforcer *rbac.Enforcer, authEventTopic string) *AuthLogic {
	l := &AuthLogic{
		BaseLogic: baseLogic,
		auth:      auth.NewAuth(baseLogic.repos, baseLogic.redisCli, jwtConfig, enforcer, authEventTopic),
	}
	return l
}

func (l *AuthLogic) AuthCheck(ctx context.Context, input dto.AuthCheckRequestDto) (result *dto.AuthCheckResponseDto, err error) {
	result, err = l.auth.AuthCheck(ctx, input)
	if err != nil {
		return
	}
	return result, nil
}

func (l *AuthLogic) RebuildPolicy(ctx context.Context, organizationType constants.UserType) (err error) {
	err = l.auth.RebuildPolicy(ctx, organizationType)

	return nil
}

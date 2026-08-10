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

// NewAuthLogic creates a new authentication logic instance.
// jwtConfig comes from the embedded BaseLogic (single shared key).
// subscribeTopics are the redis channels this replica listens on for reload
// events; publishTopics routes a userType to its broadcast channel.
func NewAuthLogic(baseLogic *BaseLogic, enforcer *rbac.Enforcer, subscribeTopics []string, publishTopics map[constants.UserType]string) *AuthLogic {
	l := &AuthLogic{
		BaseLogic: baseLogic,
		auth:      auth.NewAuth(baseLogic.repos, baseLogic.redisCli, baseLogic.jwtConfig, enforcer, auth.AuthTopics{Subscribe: subscribeTopics, Publish: publishTopics}),
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

	return
}

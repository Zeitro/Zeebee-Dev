package logic

import (
	"context"

	"zeebee-backend/common/rbac"
	"zeebee-backend/common/redisx"
	"zeebee-backend/domains/zeebee/auth"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/repo"
)

type AuthLogic struct {
	repos *repo.AllRepos
	auth  *auth.Auth
}

func NewAuthLogic(repos *repo.AllRepos, redisCli *redisx.Redisx, jwtConfig dto.JwtConfig, enforcer *rbac.Enforcer, authEventTopic string) *AuthLogic {
	return &AuthLogic{
		repos: repos,
		auth:  auth.NewAuth(repos, redisCli, jwtConfig, enforcer, authEventTopic),
	}
}

func (l *AuthLogic) AuthCheck(ctx context.Context, input dto.AuthCheckRequestDto) (result *dto.AuthCheckResponseDto, err error) {

	result, err = l.auth.AuthCheck(ctx, input)
	if err != nil {
		return
	}
	return result, nil
}

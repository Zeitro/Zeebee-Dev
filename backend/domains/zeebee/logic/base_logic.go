package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/redisx"
	"zeebee-backend/domains/zeebee/auth"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/repo"
)

// BaseLogic provides base business logic for all logic types
type BaseLogic struct {
	repos      *repo.AllRepos
	redisCli   *redisx.Client
	logger     logx.Logger
	jwtConfigs map[constants.UserType]dto.JwtConfig
}

// NewBaseLogic creates a new base logic instance.
// jwtConfigs maps each UserType the serving process handles to its JwtConfig:
// a single-platform backend passes one entry; tenant (multi-platform) passes all.
func NewBaseLogic(repos *repo.AllRepos, redisCli *redisx.Client, logger logx.Logger, jwtConfigs map[constants.UserType]dto.JwtConfig) *BaseLogic {
	return &BaseLogic{
		repos:      repos,
		redisCli:   redisCli,
		logger:     logger,
		jwtConfigs: jwtConfigs,
	}
}

// JwtConfigFor returns the JwtConfig for the given userType; missing = config error.
// Token issuance / claim resolution resolve the platform key here, so callers do
// not thread jwtConfig through every signature.
func (l *BaseLogic) JwtConfigFor(userType constants.UserType) (dto.JwtConfig, error) {
	c, ok := l.jwtConfigs[userType]
	if !ok {
		return dto.JwtConfig{}, errorx.New("jwt config not configured for user type: " + string(userType))
	}
	return c, nil
}

func (l *BaseLogic) publishAuthEvent(ctx context.Context, topic string, event auth.Msg) (err error) {
	defer func() {
		if err != nil {
			logx.WithContext(ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("topic", topic),
				logx.Field("event", event),
			).Error("failed to publish auth event")
		}
	}()

	if l.redisCli == nil {
		err = errorx.New("redis client is nil")
		return
	}

	if len(topic) == 0 {
		err = errorx.New("topic is empty")
		return
	}
	err = l.redisCli.GetClient().Publish(ctx, topic, event.String()).Err()
	return errorx.WithStack(err)
}

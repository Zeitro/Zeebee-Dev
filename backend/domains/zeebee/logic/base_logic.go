package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/redisx"
	"zeebee-backend/domains/zeebee/auth"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/repo"
)

// BaseLogic provides base business logic for all logic types
type BaseLogic struct {
	repos     *repo.AllRepos
	redisCli  *redisx.Client
	logger    logx.Logger
	jwtConfig dto.JwtConfig
}

// NewBaseLogic creates a new base logic instance.
// jwtConfig is the single shared JWT config: userType is carried in the token
// payload, so signing / verification no longer key off userType.
func NewBaseLogic(repos *repo.AllRepos, redisCli *redisx.Client, logger logx.Logger, jwtConfig dto.JwtConfig) *BaseLogic {
	return &BaseLogic{
		repos:     repos,
		redisCli:  redisCli,
		logger:    logger,
		jwtConfig: jwtConfig,
	}
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

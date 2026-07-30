package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/redisx"
	"zeebee-backend/domains/zeebee/auth"
	"zeebee-backend/domains/zeebee/repo"
)

// BaseLogic provides base business logic for all logic types
type BaseLogic struct {
	repos    *repo.AllRepos
	redisCli *redisx.Client
	logger   logx.Logger
}

// NewBaseLogic creates a new base logic instance
func NewBaseLogic(repos *repo.AllRepos, redisCli *redisx.Client, logger logx.Logger) *BaseLogic {
	return &BaseLogic{
		repos:    repos,
		redisCli: redisCli,
		logger:   logger,
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

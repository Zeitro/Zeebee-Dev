package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/redisx"
	"zeebee-backend/domains/zeebee/auth"
	"zeebee-backend/domains/zeebee/repo"
)

type baseLogic struct {
	repos *repo.AllRepos
}

func (l *baseLogic) publishAuthEvent(ctx context.Context, redisCli *redisx.Redisx, topic string, event auth.Msg) (err error) {
	defer func() {
		if err != nil {
			logx.WithContext(ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("topic", topic),
				logx.Field("event", event),
			).Error("Failed to publish auth event")
		}
	}()

	if redisCli == nil {
		err = errorx.New("Redis client is nil")
		return
	}

	if len(topic) == 0 {
		err = errorx.New("Topic is empty")
		return
	}
	err = redisCli.GetClient().Publish(ctx, topic, event.String()).Err()
	return errorx.WithStack(err)
}

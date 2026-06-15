package organizationsubscriptionlogic

import (
	"context"
	"time"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"

	"github.com/zeromicro/go-zero/core/logx"
	"zeebee-backend/common/errorx"
)

type UpdateExpireDateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewUpdateExpireDateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateExpireDateLogic {
	return &UpdateExpireDateLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *UpdateExpireDateLogic) UpdateExpireDate(in *admin.OrganizationSubscriptionUpdateExpireDateRequest) (resp *admin.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()

	if in.Data.Id <= 0 {
		return nil, errorx.New("invalid subscription id")
	}

	expireDate, err := time.Parse(time.RFC3339, in.Data.ExpireDate)
	if err != nil {
		return nil, errorx.New("invalid expire date format, expected RFC3339")
	}

	err = l.baseLogic.OrganizationSubscriptionLogic.UpdateExpireDate(l.ctx, in.Data.Id, expireDate)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &admin.EmptyResponse{}, nil
}

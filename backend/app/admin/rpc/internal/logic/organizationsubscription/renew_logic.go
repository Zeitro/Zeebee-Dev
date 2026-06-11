package organizationsubscriptionlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/data/dto"

	"github.com/zeromicro/go-zero/core/logx"
)

type RenewLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewRenewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RenewLogic {
	return &RenewLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *RenewLogic) Renew(in *admin.OrganizationSubscriptionRenewRequest) (resp *admin.OrganizationSubscriptionRenewResponse, err error) {

	dtoData, err := copierx.Trans[admin.OrganizationSubscriptionRenewRequestData, dto.OrganizationSubscriptionRenewDto](in.Data, copierx.StrToSqlTimePtr(timex.LayoutUsStyle))
	if err != nil {
		return nil, err
	}

	result, err := l.baseLogic.OrganizationSubscriptionLogic.Renew(l.ctx, in.Data.Id, *dtoData)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	protoData, err := copierx.Trans[dto.OrganizationSubscriptionEntity, admin.OrganizationSubscriptionEntity](result, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &admin.OrganizationSubscriptionRenewResponse{Data: protoData}, nil
}

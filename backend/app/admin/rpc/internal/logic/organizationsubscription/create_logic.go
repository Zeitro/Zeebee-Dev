package organizationsubscriptionlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *CreateLogic) Create(in *admin.OrganizationSubscriptionCreateRequest) (resp *admin.OrganizationSubscriptionCreateResponse, err error) {

	// organization_id comes from API layer context, passed via request metadata
	dtoData, err := copierx.Trans[admin.OrganizationSubscriptionCreateRequestData, dto.OrganizationSubscriptionCreateDto](in.Data)
	if err != nil {
		return nil, err
	}

	result, err := l.baseLogic.OrganizationSubscriptionLogic.Create(l.ctx, *dtoData)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	protoData, err := copierx.Trans[dto.OrganizationSubscriptionEntity, admin.OrganizationSubscriptionEntity](result, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &admin.OrganizationSubscriptionCreateResponse{Data: protoData}, nil
}

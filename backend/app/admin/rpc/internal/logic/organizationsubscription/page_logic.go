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
)

type PageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageLogic {
	return &PageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *PageLogic) Page(in *admin.OrganizationSubscriptionPageRequest) (resp *admin.OrganizationSubscriptionPageResponse, err error) {

	req, err := copierx.Trans[admin.OrganizationSubscriptionPageRequestData, dto.OrganizationSubscriptionPageRequestDto](in.Data)
	if err != nil {
		return nil, err
	}

	pageResult, err := l.baseLogic.OrganizationSubscriptionLogic.Page(l.ctx, *req)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	pageDto := dto.PageResultDto[dto.OrganizationSubscriptionPageEntity](pageResult)
	protoData, err := copierx.Trans[dto.PageResultDto[dto.OrganizationSubscriptionPageEntity], admin.OrganizationSubscriptionPageResponseData](&pageDto, copierx.TimeToRFC3339Str(), copierx.SqlTimePtrToStr(timex.LayoutUsStyle))
	if err != nil {
		return nil, err
	}

	return &admin.OrganizationSubscriptionPageResponse{Data: protoData}, nil
}

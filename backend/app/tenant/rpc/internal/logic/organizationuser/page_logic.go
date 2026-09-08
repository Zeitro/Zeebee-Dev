package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
)

type PageLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageLogic {
	return &PageLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

// Page returns a paginated list of organization members.
func (l *PageLogic) Page(in *tenant.OrganizationUserPageRequest) (resp *tenant.OrganizationUserPageResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return
	}

	pageReq, err := copierx.Trans[tenant.OrganizationUserPageRequestData, dto.OrganizationUserPageRequestDto](in.Data)
	if err != nil {
		return
	}
	pageReq.Filter.UserType = userType

	page, err := l.baseLogic.OrganizationUserLogic.Page(l.ctx, *pageReq)
	if err != nil {
		return
	}

	data, err := copierx.Trans[dto.PageResultDto[*dto.OrganizationUserEntityDto], tenant.OrganizationUserPageResponseData](&page, copierx.TimeToRFC3339Str())
	if err != nil {
		return
	}
	return &tenant.OrganizationUserPageResponse{Data: data}, nil
}

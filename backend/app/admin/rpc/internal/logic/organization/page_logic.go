package organizationlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
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

func (l *PageLogic) Page(in *admin.OrganizationPageRequest) (resp *admin.OrganizationPageResponse, err error) {

	pageReq, err := copierx.Trans[admin.OrganizationPageRequestData, dto.OrganizationPageRequestDto](in.Data)
	if err != nil {
		return nil, err
	}

	pageResult, err := l.baseLogic.OrganizationLogic.Page(l.ctx, *pageReq)
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[dto.PageResultDto[dto.OrganizationPageEntityDto], admin.OrganizationPageResponseData](&pageResult, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &admin.OrganizationPageResponse{
		Data: data,
	}, nil
}

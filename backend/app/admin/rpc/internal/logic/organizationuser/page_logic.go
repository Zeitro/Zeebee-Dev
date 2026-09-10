package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"
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

func (l *PageLogic) Page(in *admin.OrganizationUserPageRequest) (resp *admin.OrganizationUserPageResponse, err error) {

	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if in.Data == nil {
		return nil, logic.ErrorInvalidRequest
	}
	request := &tenant.OrganizationUserPageRequest{Data: &tenant.OrganizationUserPageRequestData{
		Page:     in.Data.Page,
		Limit:    in.Data.Limit,
		UserType: string(claim.UserType),
		Filter: &tenant.OrganizationUserPageRequestFilter{
			OrganizationId: claim.OrganizationId,
			UserId:         claim.UserId,
		},
	}}
	if in.Data.Filter != nil {
		request.Data.Filter.Keyword = in.Data.Filter.Keyword
	}
	if in.Data.Sorter != nil {
		request.Data.Sorter = &proto.CommonPageSorter{Column: in.Data.Sorter.Column, Sort: in.Data.Sorter.Sort}
	}
	results, err := l.svcCtx.TenantClient.OrganizationUser.Page(l.ctx, request)
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[tenant.OrganizationUserPageResponseData, admin.OrganizationUserPageResponseData](results.Data, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &admin.OrganizationUserPageResponse{
		Data: data,
	}, nil
}

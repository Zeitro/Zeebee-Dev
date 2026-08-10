package rolelogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
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

func (l *PageLogic) Page(in *proto.CommonPageRequest) (resp *lender.RolePageResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	pageReq, err := copierx.Trans[proto.CommonPageRequestData, tenant.RolePageRequestData](in.Data)
	if err != nil {
		return nil, err
	}
	pageReq.UserType = string(constants.UserTypeLender)
	pageReq.OrganizationId = claim.OrganizationId

	result, err := l.svcCtx.TenantClient.Role.Page(l.ctx, &tenant.RolePageRequest{Data: pageReq})
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[tenant.RolePageResponseData, lender.RolePageResponseData](result.Data, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &lender.RolePageResponse{
		Data: data,
	}, nil
}

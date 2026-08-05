package rolelogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
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

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx, constants.UserTypeLender)
	if err != nil {
		return nil, err
	}

	pageReq, err := copierx.Trans[proto.CommonPageRequestData, dto.RolePageRequestDto](in.Data)
	if err != nil {
		return nil, err
	}
	pageReq.Filter.RoleType = constants.UserTypeLender
	pageReq.Filter.OrganizationId = claim.OrganizationId
	// pageReq.Filter.UserId = claim.OrganizationUserId

	result, err := l.baseLogic.RoleLogic.Page(l.ctx, *pageReq)
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[gormx.PageResponse[dto.RoleEntityDto], lender.RolePageResponseData](&result, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &lender.RolePageResponse{
		Data: data,
	}, nil
}

package rolelogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
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

func (l *PageLogic) Page(in *proto.CommonPageRequest) (resp *borrower.RolePageResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	pageReq := dto.RolePageRequestDto{}
	_ = copierx.Assign[proto.CommonPageRequestData, dto.RolePageRequestDto](in.Data, &pageReq)

	pageReq.Filter.RoleType = constants.UserTypeBorrower
	pageReq.Filter.OrganizationId = claim.OrganizationId

	result, err := l.baseLogic.RoleLogic.Page(l.ctx, pageReq)
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[gormx.PageResponse[dto.RoleEntityDto], borrower.RolePageResponseData](&result, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &borrower.RolePageResponse{
		Data: data,
	}, nil
}

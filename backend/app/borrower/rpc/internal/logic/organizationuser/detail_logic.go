package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
)

type DetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *DetailLogic) Detail(in *proto.PathIdRequest) (resp *borrower.OrganizationUserDetailResponse, err error) {

	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	organizationUser, err := l.svcCtx.TenantClient.OrganizationUser.Detail(l.ctx, &tenant.OrganizationUserDetailRequest{Data: &tenant.OrganizationUserDetailRequestData{
		UserType:       string(constants.UserTypeBorrower),
		OrganizationId: claim.OrganizationId,
		Id:             in.Data.Id,
	}})
	if err != nil {
		return nil, err
	}

	return &borrower.OrganizationUserDetailResponse{
		Data: organizationUser.Data,
	}, nil
}

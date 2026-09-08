package organizationemaillogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
)

type SaveEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewSaveEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveEmailLogic {
	return &SaveEmailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *SaveEmailLogic) SaveEmail(in *borrower.OrganizationEmailSaveEmailRequest) (resp *proto.EmptyResponse, err error) {

	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	input, err := copierx.Trans[borrower.OrganizationEmailSaveEmailRequestData, tenant.OrganizationEmailSaveEmailRequestData](in.Data)
	if err != nil {
		return nil, err
	}
	input.OrganizationId = claim.OrganizationId
	input.UserType = string(constants.UserTypeBorrower)
	if input.AuthStrategy == string(constants.EmailAuthStrategyUser) {
		input.UserId = claim.UserId
	}

	_, err = l.svcCtx.TenantClient.OrganizationEmail.SaveEmail(l.ctx, &tenant.OrganizationEmailSaveEmailRequest{Data: input})
	if err != nil {
		return nil, err
	}

	return &proto.EmptyResponse{}, nil
}

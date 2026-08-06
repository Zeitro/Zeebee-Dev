package organizationemaillogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
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

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	input, err := copierx.Trans[borrower.OrganizationEmailSaveEmailRequestData, dto.OrganizationEmailSaveEmailRequestDto](in.Data)
	if err != nil {
		return nil, err
	}
	input.OrganizationId = claim.OrganizationId
	input.OrganizationType = constants.UserTypeBorrower
	if input.AuthStrategy == constants.EmailAuthStrategyUser {
		input.UserId = claim.UserId
	}

	err = l.baseLogic.OrganizationEmailLogic.SaveEmail(l.ctx, *input)
	if err != nil {
		return nil, err
	}

	return &proto.EmptyResponse{}, nil
}

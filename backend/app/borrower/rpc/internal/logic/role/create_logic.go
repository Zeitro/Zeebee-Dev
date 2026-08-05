package rolelogic

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

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *CreateLogic) Create(in *borrower.RoleCreateRequest) (resp *borrower.RoleCreateResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}
	input := &dto.RoleCreateRequestDto{}
	input, err = copierx.Trans[borrower.RoleCreateRequestData, dto.RoleCreateRequestDto](in.Data)
	if err != nil {
		return nil, err
	}

	input.Type = constants.UserTypeBorrower
	input.OrganizationId = claim.OrganizationId

	role, err := l.baseLogic.RoleLogic.Create(l.ctx, *input)
	if err != nil {
		return nil, err
	}

	respData := &proto.RoleEntity{}
	respData, err = copierx.Trans[dto.RoleEntityDto, proto.RoleEntity](role, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &borrower.RoleCreateResponse{
		Data: respData,
	}, nil
}

package rolelogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
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

func (l *CreateLogic) Create(in *admin.RoleCreateRequest) (resp *admin.RoleDetailResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx, constants.UserTypeAdmin)
	if err != nil {
		return nil, err
	}

	input, err := copierx.Trans[admin.RoleCreateRequestData, dto.RoleCreateRequestDto](in.Data)
	if err != nil {
		return nil, err
	}

	input.Type = constants.UserTypeAdmin
	input.OrganizationId = claim.OrganizationId

	role, err := l.baseLogic.RoleLogic.Create(l.ctx, *input)
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[dto.RoleEntityDto, proto.RoleEntity](role, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	return &admin.RoleDetailResponse{
		Data: data,
	}, nil
}

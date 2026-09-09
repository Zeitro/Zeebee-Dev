package userlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
)

type CreateAdminUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewCreateAdminUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAdminUserLogic {
	return &CreateAdminUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// CreateAdminUser provisions an admin-backend user (ODD-4): super-admin creates
// a user attached to an existing organization with the given roles. Always admin
// user_type regardless of request.
func (l *CreateAdminUserLogic) CreateAdminUser(in *tenant.CreateAdminUserRequest) (resp *tenant.CreateAdminUserResponse, err error) {
	if in.Data == nil {
		return nil, errorx.New("request data is nil")
	}

	input, err := copierx.Trans[tenant.CreateAdminUserRequestData, dto.OrganizationUserCreateAdminRequestDto](in.Data)
	if err != nil {
		return nil, err
	}
	input.UserType = constants.UserTypeAdmin

	orgUser, err := l.baseLogic.OrganizationUserLogic.CreateAdminOrganizationUser(l.ctx, *input)
	if err != nil {
		return nil, err
	}

	return &tenant.CreateAdminUserResponse{
		Data: &tenant.CreateAdminUserResponseData{
			UserId: orgUser.UserId,
		},
	}, nil
}

package resourcelogic

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

type UserMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewUserMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserMenuLogic {
	return &UserMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *UserMenuLogic) UserMenu(in *proto.EmptyRequest) (resp *admin.ResourceTreeResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx, constants.UserTypeAdmin)
	if err != nil {
		return nil, err
	}
	input := dto.ResourceUserMenuRequestDto{
		OrganizationId:     claim.OrganizationId,
		OrganizationUserId: claim.OrganizationUserId,
		Type:               claim.UserType,
	}

	result, err := l.baseLogic.ResourceLogic.GetUserResources(l.ctx, input)
	if err != nil {
		return nil, err
	}
	// trans response
	resp = &admin.ResourceTreeResponse{}
	menus, err := copierx.Trans[[]*dto.ResourceTreeDto, []*proto.ResourceTreeEntity](&result, copierx.TimeToRFC3339Str())

	return &admin.ResourceTreeResponse{
		Data: *menus,
	}, nil
}

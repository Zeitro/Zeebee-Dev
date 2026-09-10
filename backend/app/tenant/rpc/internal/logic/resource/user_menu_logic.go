package resourcelogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
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
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// UserMenu derives the user's menu tree from roles + subscription gate.
func (l *UserMenuLogic) UserMenu(in *tenant.UserMenuRequest) (resp *tenant.ResourceTreeResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return nil, err
	}

	input := dto.ResourceUserMenuRequestDto{
		OrganizationId:     in.Data.OrganizationId,
		OrganizationUserId: in.Data.OrganizationUserId,
		Type:               userType,
	}
	result, err := l.baseLogic.ResourceLogic.GetUserResources(l.ctx, input)
	if err != nil {
		return nil, err
	}

	menus, err := copierx.Trans[[]*dto.ResourceTreeDto, []*proto.ResourceTreeEntity](&result, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	return &tenant.ResourceTreeResponse{Data: *menus}, nil
}

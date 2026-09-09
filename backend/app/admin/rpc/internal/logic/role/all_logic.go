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

type AllLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllLogic {
	return &AllLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *AllLogic) All(in *admin.RoleAllRequest) (resp *admin.RoleAllResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	result, err := l.baseLogic.RoleLogic.All(l.ctx, claim.OrganizationId, constants.UserTypeAdmin)
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[[]*dto.RoleEntityDto, []*proto.RoleEntity](&result, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &admin.RoleAllResponse{
		Data: *data,
	}, nil
}

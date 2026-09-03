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

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// Update edits a resource.
func (l *UpdateLogic) Update(in *tenant.ResourceUpdateRequest) (resp *tenant.ResourceDetailResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	if _, err = logic.ParseUserType(in.Data.UserType); err != nil {
		return nil, err
	}

	input, err := copierx.Trans[tenant.ResourceUpdateRequestData, dto.ResourceUpdateRequestDto](in.Data)
	if err != nil {
		return nil, err
	}

	result, err := l.baseLogic.ResourceLogic.UpdateResource(l.ctx, input)
	if err != nil {
		return nil, err
	}

	respData, err := copierx.Trans[dto.ResourceEntityDto, proto.ResourceEntity](result, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	return &tenant.ResourceDetailResponse{Data: respData}, nil
}

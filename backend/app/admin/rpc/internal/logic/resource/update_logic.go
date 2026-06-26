package resourcelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
)

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *UpdateLogic) Update(in *admin.ResourceUpdateRequest) (resp *admin.ResourceDetailResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()

	data := in.Data
	if data == nil {
		return nil, logic.ErrorInvalidRequest
	}
	input, err := copierx.Trans[admin.ResourceUpdateRequestData, dto.ResourceUpdateRequestDto](data)
	if err != nil {
		return nil, err
	}

	result, err := l.baseLogic.ResourceLogic.UpdateResource(l.ctx, input)
	if err != nil {
		return nil, err
	}

	respData, err := copierx.Trans[dto.ResourceEntityDto, admin.ResourceEntity](result, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &admin.ResourceDetailResponse{
		Data: respData,
	}, nil
}

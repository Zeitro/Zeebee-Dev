package apilogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
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
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *UpdateLogic) Update(in *admin.ApiUpdateRequest) (resp *admin.ApiDetailResponse, err error) {
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Error(err)
		}
	}()

	data := in.Data
	if data == nil {
		return nil, logic.ErrorInvalidRequest
	}
	input, err := copierx.Trans[admin.ApiUpdateRequestData, dto.ApiUpdateRequestDto](data)
	if err != nil {
		return nil, err
	}
	result, err := l.baseLogic.ApiLogic.UpdateApi(l.ctx, input)
	if err != nil {
		return nil, err
	}
	respData, err := copierx.Trans[dto.ApiEntityDto, proto.ApiEntity](result, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &admin.ApiDetailResponse{
		Data: respData,
	}, nil
}

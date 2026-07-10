package basetemplatelogic

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

// Update updates an existing base template.
// Tag field is immutable and cannot be modified (FR-016).
func (l *UpdateLogic) Update(in *admin.UpdateBaseTemplateRequest) (resp *admin.UpdateBaseTemplateResponse, err error) {
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Error(err)
		}
	}()

	// Convert proto request to DTO
	updateReq, err := copierx.Trans[admin.UpdateBaseTemplateRequestData, dto.BaseTemplateUpdateDto](in.Data)
	if err != nil {
		return nil, err
	}

	// Perform update
	result, err := l.baseLogic.BaseTemplateLogic.Update(l.ctx, *updateReq)
	if err != nil {
		return nil, err
	}

	// Convert BO to proto for response
	protoData, err := copierx.Trans[dto.BaseTemplateEntityDto, proto.BaseTemplateEntity](result, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &admin.UpdateBaseTemplateResponse{
		Data: protoData,
	}, nil
}

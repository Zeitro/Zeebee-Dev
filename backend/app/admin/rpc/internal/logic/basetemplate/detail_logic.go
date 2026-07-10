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

type DetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// Detail retrieves base template details by ID
func (l *DetailLogic) Detail(in *admin.GetBaseTemplateRequest) (resp *admin.GetBaseTemplateResponse, err error) {
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Error(err)
		}
	}()

	// Get base template details from logic layer
	result, err := l.baseLogic.BaseTemplateLogic.Detail(l.ctx, in.Data.Data)
	if err != nil {
		return nil, err
	}

	// Convert BO to proto for response
	protoData, err := copierx.Trans[dto.BaseTemplateEntityDto, proto.BaseTemplateEntity](result, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &admin.GetBaseTemplateResponse{
		Data: protoData,
	}, nil
}

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

// Create creates a new base template with validation
func (l *CreateLogic) Create(in *admin.CreateBaseTemplateRequest) (resp *admin.CreateBaseTemplateResponse, err error) {
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Error(err)
		}
	}()

	// Convert proto request to DTO
	dtoData, err := copierx.Trans[admin.CreateBaseTemplateRequestData, dto.BaseTemplateCreateDto](in.Data)
	if err != nil {
		return nil, err
	}

	// Create base template
	result, err := l.baseLogic.BaseTemplateLogic.Create(l.ctx, *dtoData)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// Convert BO to proto for response
	protoData, err := copierx.Trans[dto.BaseTemplateEntityDto, proto.BaseTemplateEntity](result, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &admin.CreateBaseTemplateResponse{
		Data: protoData,
	}, nil
}

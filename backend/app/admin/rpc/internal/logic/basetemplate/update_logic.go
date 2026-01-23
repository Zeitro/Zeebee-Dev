package basetemplatelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/bo"
	"zeebee-backend/domains/zeebee/data/dto"
)

type UpdateLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

// Update updates an existing base template.
// Tag field is immutable and cannot be modified (FR-016).
func (l *UpdateLogic) Update(in *admin.UpdateBaseTemplateRequest) (resp *admin.UpdateBaseTemplateResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	// Get existing template to verify immutability
	existing, err := l.baseLogic.BaseTemplateLogic.GetByIDForUpdate(l.ctx, in.Data.Id)
	if err != nil {
		return nil, err
	}

	// Validate tag immutability (FR-016: tag cannot be changed after creation)
	if existing.Tag != in.Data.Tag {
		return nil, errorx.New("Template tag is immutable and cannot be changed")
	}

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
	protoData, err := copierx.Trans[bo.BaseTemplateBo, admin.BaseTemplateEntity](result, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &admin.UpdateBaseTemplateResponse{
		Data: protoData,
	}, nil
}

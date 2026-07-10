package basetemplatelogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/data/dto"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type PreviewLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewPreviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreviewLogic {
	return &PreviewLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *PreviewLogic) Preview(in *admin.BaseTemplatePreviewRequest) (resp *admin.BaseTemplatePreviewResponse, err error) {
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Error(err)
		}
	}()
	// Convert proto request to DTO
	req, err := copierx.Trans[admin.BaseTemplatePreviewRequestData, dto.BaseTemplatePreviewDto](in.Data)
	if err != nil {
		return nil, err
	}
	previewResult, err := l.baseLogic.BaseTemplateLogic.Preview(l.ctx, *req)
	if err != nil {
		return nil, err
	}

	return &admin.BaseTemplatePreviewResponse{
		Data: &admin.BaseTemplatePreviewResponseData{
			Content: previewResult,
		},
	}, nil
}

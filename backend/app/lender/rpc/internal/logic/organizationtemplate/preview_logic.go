package organizationtemplatelogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/data/dto"
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

func (l *PreviewLogic) Preview(in *lender.OrganizationTemplatePreviewRequest) (resp *lender.OrganizationTemplatePreviewResponse, err error) {

	input, err := copierx.Trans[lender.OrganizationTemplatePreviewRequestData, dto.BaseTemplatePreviewDto](in.Data)

	result, err := l.baseLogic.BaseTemplateLogic.Preview(l.ctx, *input)

	return &lender.OrganizationTemplatePreviewResponse{
		Data: &lender.OrganizationTemplatePreviewResponseData{
			Content: result,
		},
	}, err
}

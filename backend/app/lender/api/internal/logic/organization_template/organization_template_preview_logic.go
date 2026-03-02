package organization_template

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationTemplatePreviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Preview
func NewOrganizationTemplatePreviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationTemplatePreviewLogic {
	return &OrganizationTemplatePreviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationTemplatePreviewLogic) OrganizationTemplatePreview(req *types.OrganizationTemplatePreviewRequest) (resp *types.OrganizationTemplatePreviewResponse, err error) {
	rpcReq := &lender.OrganizationTemplatePreviewRequest{}

	rpcReq.Data, err = copierx.Trans[types.OrganizationTemplatePreviewRequest, lender.OrganizationTemplatePreviewRequestData](req)
	if err != nil {
		return
	}

	rpcResp, err := l.svcCtx.OrganizationTemplateRpcClient.Preview(l.ctx, rpcReq)
	if err != nil {
		return
	}
	resp, err = copierx.Trans[lender.OrganizationTemplatePreviewResponseData, types.OrganizationTemplatePreviewResponse](rpcResp.Data)

	return resp, err
}

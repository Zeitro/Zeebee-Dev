package organization_template

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationTemplatePageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// page
func NewOrganizationTemplatePageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationTemplatePageLogic {
	return &OrganizationTemplatePageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationTemplatePageLogic) OrganizationTemplatePage(req *types.OrganizationTemplatePageRequest) (resp *types.OrganizationTemplatePageResponse, err error) {
	rpcReq := &lender.OrganizationTemplatePageRequest{}

	rpcReq.Data, err = copierx.Trans[types.OrganizationTemplatePageRequest, lender.OrganizationTemplatePageRequestData](req)
	if err != nil {
		return
	}

	rpcResp, err := l.svcCtx.OrganizationTemplateRpcClient.Page(l.ctx, rpcReq)
	if err != nil {
		return
	}
	resp, err = copierx.Trans[lender.OrganizationTemplatePageResponseData, types.OrganizationTemplatePageResponse](rpcResp.Data)

	return resp, err
}

package organization_template

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationTemplateUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Update
func NewOrganizationTemplateUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationTemplateUpdateLogic {
	return &OrganizationTemplateUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationTemplateUpdateLogic) OrganizationTemplateUpdate(req *types.OrganizationTemplateUpdateRequest) (resp *types.OrganizationTemplateDetailResponse, err error) {
	rpcReq := &lender.UpdateOrganizationTemplateRequest{}

	rpcReq.Data, err = copierx.Trans[types.OrganizationTemplateUpdateRequest, lender.UpdateOrganizationTemplateRequestData](req)
	if err != nil {
		return
	}

	rpcResp, err := l.svcCtx.OrganizationTemplateRpcClient.Update(l.ctx, rpcReq)
	if err != nil {
		return
	}
	resp, err = copierx.Trans[proto.OrganizationTemplateEntity, types.OrganizationTemplateDetailResponse](rpcResp.Data)

	return resp, err
}

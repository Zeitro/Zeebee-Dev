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

type OrganizationTemplateDetailByPurposeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get Detail by purpose
func NewOrganizationTemplateDetailByPurposeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationTemplateDetailByPurposeLogic {
	return &OrganizationTemplateDetailByPurposeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationTemplateDetailByPurposeLogic) OrganizationTemplateDetailByPurpose(req *types.OrganizationTemplateDetailByPurposeRequest) (resp *types.OrganizationTemplateDetailResponse, err error) {
	rpcReq := &lender.OrganizationTemplateDetailByPurposeRequest{}

	rpcReq.Data, err = copierx.Trans[types.OrganizationTemplateDetailByPurposeRequest, lender.OrganizationTemplateDetailByPurposeRequestData](req)
	if err != nil {
		return
	}

	rpcResp, err := l.svcCtx.OrganizationTemplateRpcClient.DetailByPurpose(l.ctx, rpcReq)
	if err != nil {
		return
	}
	resp, err = copierx.Trans[proto.OrganizationTemplateEntity, types.OrganizationTemplateDetailResponse](rpcResp.Data)

	return resp, err
}

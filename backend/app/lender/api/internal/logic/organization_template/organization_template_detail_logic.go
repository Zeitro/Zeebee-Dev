package organization_template

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationTemplateDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Detail
func NewOrganizationTemplateDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationTemplateDetailLogic {
	return &OrganizationTemplateDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationTemplateDetailLogic) OrganizationTemplateDetail(req *types.PathIdRequest) (resp *types.OrganizationTemplateDetailResponse, err error) {
	rpcReq := &lender.IdRequest{Data: req.Id}

	rpcResp, err := l.svcCtx.OrganizationTemplateRpcClient.Detail(l.ctx, rpcReq)
	if err != nil {
		return
	}
	resp, err = copierx.Trans[lender.OrganizationTemplateEntity, types.OrganizationTemplateDetailResponse](rpcResp.Data)

	return resp, err
}

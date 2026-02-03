package organization

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// page
func NewOrganizationPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationPageLogic {
	return &OrganizationPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationPageLogic) OrganizationPage(req *types.OrganizationPageRequest) (resp *types.OrganizationPageResponse, err error) {
	rpcReq := &admin.OrganizationPageRequest{}
	rpcReq.Data, err = copierx.Trans[types.OrganizationPageRequest, admin.OrganizationPageRequestData](req)

	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.OrganizationRpcClient.Page(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	resp, err = copierx.Trans[admin.OrganizationPageResponseData, types.OrganizationPageResponse](rpcResp.Data)

	return resp, err
}

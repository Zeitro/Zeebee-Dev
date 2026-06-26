package organization_subscription

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationSubscriptionPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// page
func NewOrganizationSubscriptionPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationSubscriptionPageLogic {
	return &OrganizationSubscriptionPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationSubscriptionPageLogic) OrganizationSubscriptionPage(req *types.OrganizationSubscriptionPageRequest) (resp *types.OrganizationSubscriptionPageResponse, err error) {

	// Convert types to pb request
	rpcReq := &admin.OrganizationSubscriptionPageRequest{}
	rpcReq.Data, err = copierx.Trans[types.OrganizationSubscriptionPageRequest, admin.OrganizationSubscriptionPageRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	// Call RPC server
	rpcResp, err := l.svcCtx.OrganizationSubscriptionRpcClient.Page(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// Convert pb response to types
	resp, err = copierx.Trans[admin.OrganizationSubscriptionPageResponseData, types.OrganizationSubscriptionPageResponse](rpcResp.Data)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return resp, nil
}

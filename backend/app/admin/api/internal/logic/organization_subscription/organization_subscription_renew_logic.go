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

type OrganizationSubscriptionRenewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// renew
func NewOrganizationSubscriptionRenewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationSubscriptionRenewLogic {
	return &OrganizationSubscriptionRenewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationSubscriptionRenewLogic) OrganizationSubscriptionRenew(req *types.OrganizationSubscriptionRenewRequest) (resp *types.OrganizationSubscriptionDetailResponse, err error) {

	// Convert types to pb request
	rpcReq := &admin.OrganizationSubscriptionRenewRequest{}
	rpcReq.Data, err = copierx.Trans[types.OrganizationSubscriptionRenewRequest, admin.OrganizationSubscriptionRenewRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	rpcReq.Data.Id = req.Id

	// Call RPC server
	rpcResp, err := l.svcCtx.OrganizationSubscriptionRpcClient.Renew(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// Convert pb response to types
	resp, err = copierx.Trans[admin.OrganizationSubscriptionEntity, types.OrganizationSubscriptionDetailResponse](rpcResp.Data, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return resp, nil
}

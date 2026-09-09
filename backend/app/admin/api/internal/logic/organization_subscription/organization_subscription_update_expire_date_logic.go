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

type OrganizationSubscriptionUpdateExpireDateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update expire date
func NewOrganizationSubscriptionUpdateExpireDateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationSubscriptionUpdateExpireDateLogic {
	return &OrganizationSubscriptionUpdateExpireDateLogic{
		Logger: svcCtx.Logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationSubscriptionUpdateExpireDateLogic) OrganizationSubscriptionUpdateExpireDate(req *types.OrganizationSubscriptionUpdateExpireDateRequest) (resp *types.EmptyResponse, err error) {
	// Convert types to pb request
	rpcReq := &admin.OrganizationSubscriptionUpdateExpireDateRequest{}
	rpcReq.Data, err = copierx.Trans[types.OrganizationSubscriptionUpdateExpireDateRequest, admin.OrganizationSubscriptionUpdateExpireDateRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	rpcReq.Data.Id = req.Id

	// Call RPC server
	_, err = l.svcCtx.OrganizationSubscriptionRpcClient.UpdateExpireDate(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return &types.EmptyResponse{}, nil
}

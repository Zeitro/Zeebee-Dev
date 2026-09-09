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

type OrganizationSubscriptionSwitchStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// switch status
func NewOrganizationSubscriptionSwitchStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationSubscriptionSwitchStatusLogic {
	return &OrganizationSubscriptionSwitchStatusLogic{
		Logger: svcCtx.Logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationSubscriptionSwitchStatusLogic) OrganizationSubscriptionSwitchStatus(req *types.OrganizationSubscriptionSwitchStatusRequest) (resp *types.EmptyResponse, err error) {

	// Convert types to pb request
	rpcReq := &admin.OrganizationSubscriptionSwitchStatusRequest{}
	rpcReq.Data, err = copierx.Trans[types.OrganizationSubscriptionSwitchStatusRequest, admin.OrganizationSubscriptionSwitchStatusRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	rpcReq.Data.Id = req.Id

	// Call RPC server
	_, err = l.svcCtx.OrganizationSubscriptionRpcClient.SwitchStatus(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &types.EmptyResponse{}, nil
}

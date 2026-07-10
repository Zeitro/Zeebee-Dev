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

type OrganizationSwitchStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// switch status
func NewOrganizationSwitchStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationSwitchStatusLogic {
	return &OrganizationSwitchStatusLogic{
		Logger: svcCtx.Logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationSwitchStatusLogic) OrganizationSwitchStatus(req *types.OrganizationSwitchStatusRequest) (resp *types.EmptyResponse, err error) {

	// Convert types to pb request
	rpcReq := &admin.OrganizationSwitchStatusRequest{
		Data: &admin.OrganizationSwitchStatusRequestData{},
	}
	rpcReq.Data, err = copierx.Trans[types.OrganizationSwitchStatusRequest, admin.OrganizationSwitchStatusRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	// Call RPC server
	_, err = l.svcCtx.OrganizationRpcClient.SwitchStatus(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &types.EmptyResponse{}, nil
}

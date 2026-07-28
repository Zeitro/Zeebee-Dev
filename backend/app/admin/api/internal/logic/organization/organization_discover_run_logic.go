package organization

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationDiscoverRunLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// lender discover run (scheduled trigger only; manual flow uses preview + commit)
func NewOrganizationDiscoverRunLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationDiscoverRunLogic {
	return &OrganizationDiscoverRunLogic{
		Logger: svcCtx.Logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationDiscoverRunLogic) LenderDiscoverRun(req *types.LenderDiscoverRunRequest) (resp *types.LenderDiscoverRunResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	rpcReq := &admin.LenderDiscoverRunRequest{
		Data: &admin.LenderDiscoverRunRequestData{Trigger: req.Trigger},
	}
	rpcResp, err := l.svcCtx.OrganizationRpcClient.DiscoverRun(context.WithoutCancel(l.ctx), rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	data, err := discoverProtoRunToTypes(rpcResp.GetData())
	if err != nil {
		return nil, err
	}
	return &types.LenderDiscoverRunResponse{Data: data}, nil
}

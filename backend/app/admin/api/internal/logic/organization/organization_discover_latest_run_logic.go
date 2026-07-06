package organization

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationDiscoverLatestRunLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrganizationDiscoverLatestRunLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationDiscoverLatestRunLogic {
	return &OrganizationDiscoverLatestRunLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationDiscoverLatestRunLogic) LenderDiscoverLatestRun(_ *types.EmptyRequest) (resp *types.LenderDiscoverRunResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	rpcResp, err := l.svcCtx.OrganizationRpcClient.DiscoverLatestRun(l.ctx, &admin.EmptyRequest{})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	data, err := discoverProtoRunToTypes(rpcResp.GetData())
	if err != nil {
		return nil, err
	}
	return &types.LenderDiscoverRunResponse{Data: data}, nil
}

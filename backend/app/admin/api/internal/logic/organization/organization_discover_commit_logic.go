package organization

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationDiscoverCommitLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrganizationDiscoverCommitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationDiscoverCommitLogic {
	return &OrganizationDiscoverCommitLogic{
		Logger: svcCtx.Logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationDiscoverCommitLogic) LenderDiscoverCommit(req *types.LenderDiscoverCommitRequest) (resp *types.LenderDiscoverRunResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	rpcReq, err := discoverTypesCommitRequestToProto(req)
	if err != nil {
		return nil, err
	}
	rpcResp, err := l.svcCtx.OrganizationRpcClient.DiscoverCommit(context.WithoutCancel(l.ctx), rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	data, err := discoverProtoRunToTypes(rpcResp.GetData())
	if err != nil {
		return nil, err
	}
	return &types.LenderDiscoverRunResponse{Data: data}, nil
}

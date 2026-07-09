package role

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"
)

type RoleAllLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// all
func NewRoleAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleAllLogic {
	return &RoleAllLogic{
		Logger: svcCtx.Logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleAllLogic) RoleAll(req *types.EmptyRequest) (resp *types.RoleAllResponse, err error) {
	resp = &types.RoleAllResponse{
		List: make([]*types.RoleEntity, 0),
	}

	rpcReq := &admin.RoleAllRequest{
		Data: &admin.RoleAllRequestData{},
	}

	rpcResp, err := l.svcCtx.RoleRpcClient.All(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// trans response
	list, err := copierx.Trans[[]*proto.RoleEntity, []*types.RoleEntity](&rpcResp.Data)
	if err != nil {
		return nil, err
	}
	resp.List = *list

	return resp, err
}

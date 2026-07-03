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

type RoleUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update
func NewRoleUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleUpdateLogic {
	return &RoleUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleUpdateLogic) RoleUpdate(req *types.RoleUpdateRequest) (resp *types.RoleEntity, err error) {
	resp = &types.RoleEntity{}
	rpcReq := &admin.RoleUpdateRequest{
		Data: &admin.RoleUpdateRequestData{},
	}

	rpcReq.Data, err = copierx.Trans[types.RoleUpdateRequest, admin.RoleUpdateRequestData](req)
	if err != nil {
		return nil, err
	}

	relations, err := copierx.Trans[[]types.RoleAndResourceRelation, []*admin.RoleAndResourceRelationEntity](&req.Resources)
	if err != nil {
		return nil, err
	}

	rpcReq.Data.RoleAndResources = *relations

	rpcResp, err := l.svcCtx.RoleRpcClient.Update(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// trans response
	resp, err = copierx.Trans[proto.RoleEntity, types.RoleEntity](rpcResp.Data)
	if err != nil {
		return nil, err
	}

	return resp, err
}

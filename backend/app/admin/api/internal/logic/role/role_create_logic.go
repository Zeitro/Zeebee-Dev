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

type RoleCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// create
func NewRoleCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleCreateLogic {
	return &RoleCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleCreateLogic) RoleCreate(req *types.RoleCreateRequest) (resp *types.RoleEntity, err error) {
	resp = &types.RoleEntity{}
	rpcReq := &admin.RoleCreateRequest{
		Data: &admin.RoleCreateRequestData{},
	}

	rpcReq.Data, err = copierx.Trans[types.RoleCreateRequest, admin.RoleCreateRequestData](req)
	if err != nil {
		return nil, err
	}

	relations, err := copierx.Trans[[]types.RoleAndResourceRelation, []*admin.RoleAndResourceRelationEntity](&req.Resources)
	if err != nil {
		return nil, err
	}

	rpcReq.Data.RoleAndResources = *relations

	rpcResp, err := l.svcCtx.RoleRpcClient.Create(l.ctx, rpcReq)
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

package role

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"
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
	err = copier.CopyWithOption(rpcReq.Data, req, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}

	relations := slice.Map[int64, *admin.RoleAndResourceRelationEntity](req.ResourceIds, func(index int, item int64) *admin.RoleAndResourceRelationEntity {
		return &admin.RoleAndResourceRelationEntity{
			RoleId:     req.Id,
			ResourceId: item,
			DataScope:  string(constants.DataScopeAll),
		}
	})
	rpcReq.Data.RoleAndResources = relations

	rpcResp, err := l.svcCtx.RoleRpcClient.Update(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// trans response
	err = copier.CopyWithOption(resp, &rpcResp.Data, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}
	return resp, err
}

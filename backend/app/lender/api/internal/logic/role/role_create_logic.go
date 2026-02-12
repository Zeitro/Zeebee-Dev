package role

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
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
	rpcReq := &lender.RoleCreateRequest{
		Data: &lender.RoleCreateRequestData{},
	}
	data, err := copierx.Trans[types.RoleCreateRequest, lender.RoleCreateRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}
	rpcReq.Data = data

	relations := slice.Map[int64, *lender.RoleAndResourceRelationEntity](req.ResourceIds, func(index int, item int64) *lender.RoleAndResourceRelationEntity {
		return &lender.RoleAndResourceRelationEntity{
			ResourceId: item,
			DataScope:  string(constants.DataScopeAll),
		}
	})
	rpcReq.Data.RoleAndResources = relations

	rpcResp, err := l.svcCtx.RoleRpcClient.Create(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// trans response
	resp, err = copierx.Trans[lender.RoleEntity, types.RoleEntity](rpcResp.Data, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}
	return resp, err
}

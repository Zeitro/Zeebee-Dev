package role

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
)

type RoleAllLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// all
func NewRoleAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleAllLogic {
	return &RoleAllLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleAllLogic) RoleAll(req *types.EmptyRequest) (resp *types.RoleAllResponse, err error) {
	rpcReq := &lender.RoleAllRequest{
		Data: &lender.RoleAllRequestData{},
	}

	rpcResp, err := l.svcCtx.RoleRpcClient.All(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// trans response
	list, err := copierx.Trans[[]*lender.RoleEntity, []*types.RoleEntity](&rpcResp.Data, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return &types.RoleAllResponse{List: *list}, nil
}

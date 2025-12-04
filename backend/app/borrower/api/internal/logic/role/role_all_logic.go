package role

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/utils"
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
	resp = &types.RoleAllResponse{
		List: make([]*types.RoleEntity, 0),
	}

	claim, err := utils.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	rpcReq := &borrower.RoleAllRequest{
		Data: &borrower.RoleAllRequestData{
			UserId:         claim.UserId,
			OrganizationId: claim.OrganizationId,
		},
	}

	rpcResp, err := l.svcCtx.RoleRpcClient.All(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// trans response
	err = copier.CopyWithOption(&resp.List, &rpcResp.Data, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return resp, nil
}

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

type RoleDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// detail
func NewRoleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleDetailLogic {
	return &RoleDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleDetailLogic) RoleDetail(req *types.PathIdRequest) (resp *types.RoleEntity, err error) {
	resp = &types.RoleEntity{}
	rpcReq := &borrower.RoleDetailRequest{
		Data: &borrower.RoleDetailRequestData{},
	}
	rpcReq.Data.Id = req.Id

	claim, err := utils.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	rpcReq.Data.UserId = claim.UserId
	rpcReq.Data.OrganizationId = claim.OrganizationId

	rpcResp, err := l.svcCtx.RoleRpcClient.Detail(l.ctx, rpcReq)
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

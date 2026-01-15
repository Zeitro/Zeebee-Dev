package organization_user

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationUserDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// delete
func NewOrganizationUserDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationUserDeleteLogic {
	return &OrganizationUserDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationUserDeleteLogic) OrganizationUserDelete(req *types.PathIdRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &lender.OrganizationUserDeleteRequest{
		Data: &lender.OrganizationUserDeleteRequestData{
			Id: req.Id,
		},
	}

	_, err = l.svcCtx.OrganizationUserRpcClient.Delete(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResponse{}, nil
}

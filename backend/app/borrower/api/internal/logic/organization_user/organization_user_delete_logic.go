package organization_user

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

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
	rpcReq := &borrower.OrganizationUserDeleteRequest{
		Data: &borrower.OrganizationUserDeleteRequestData{},
	}
	rpcReq.Data, err = copierx.Trans[types.PathIdRequest, borrower.OrganizationUserDeleteRequestData](req)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.OrganizationUserRpcClient.Delete(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	resp = &types.EmptyResponse{}
	return resp, nil
}

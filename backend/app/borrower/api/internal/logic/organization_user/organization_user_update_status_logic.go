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

type OrganizationUserUpdateStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// swaitch status
func NewOrganizationUserUpdateStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationUserUpdateStatusLogic {
	return &OrganizationUserUpdateStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationUserUpdateStatusLogic) OrganizationUserUpdateStatus(req *types.OrganizationUserSwitchStatusRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &borrower.OrganizationUserUpdateStatusRequest{
		Data: &borrower.OrganizationUserUpdateStatusRequestData{},
	}
	rpcReq.Data, err = copierx.Trans[types.OrganizationUserSwitchStatusRequest, borrower.OrganizationUserUpdateStatusRequestData](req)

	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.OrganizationUserRpcClient.SwitchStatus(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	resp = &types.EmptyResponse{}
	return resp, nil
}

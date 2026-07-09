package organization_user

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationUserSwitchStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// switch status
func NewOrganizationUserSwitchStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationUserSwitchStatusLogic {
	return &OrganizationUserSwitchStatusLogic{
		Logger: svcCtx.Logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationUserSwitchStatusLogic) OrganizationUserSwitchStatus(req *types.OrganizationUserSwitchStatusRequest) (resp *types.EmptyResponse, err error) {
	resp = &types.EmptyResponse{}
	rpcReq := &admin.OrganizationUserSwitchStatusRequest{
		Data: &admin.OrganizationUserSwitchStatusRequestData{},
	}

	rpcReq.Data, err = copierx.Trans[types.OrganizationUserSwitchStatusRequest, admin.OrganizationUserSwitchStatusRequestData](req)
	if err != nil {
		err = errorx.GetFailedTransReqErr(err)
		return nil, err
	}

	_, err = l.svcCtx.OrganizationUserRpcClient.SwitchStatus(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return resp, nil
}

package organization_user

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
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
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationUserSwitchStatusLogic) OrganizationUserSwitchStatus(req *types.OrganizationUserSwitchStatusRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &admin.OrganizationUserSwitchStatusRequest{
		Data: &admin.OrganizationUserSwitchStatusRequestData{},
	}
	err = copier.CopyWithOption(rpcReq.Data, req, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
	})
	if err != nil {
		err = errorx.GetFailedTransReqErr(err)
		return nil, err
	}

	_, err = l.svcCtx.OrganizationUserRpcClient.SwitchStatus(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	resp = &types.EmptyResponse{}
	return resp, nil
}

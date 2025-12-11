package organization_user

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationUserAssignRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// assign roles
func NewOrganizationUserAssignRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationUserAssignRolesLogic {
	return &OrganizationUserAssignRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationUserAssignRolesLogic) OrganizationUserAssignRoles(req *types.OrganizationUserAssignRolesRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &borrower.OrganizationUserAssignRolesRequest{
		Data: &borrower.OrganizationUserAssignRolesRequestData{},
	}
	err = copier.CopyWithOption(rpcReq.Data, req, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
	})
	if err != nil {
		err = errorx.GetFailedTransReqErr(err)
		return nil, err
	}

	_, err = l.svcCtx.OrganizationUserRpcClient.AssignRoles(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	resp = &types.EmptyResponse{}
	return resp, nil
}

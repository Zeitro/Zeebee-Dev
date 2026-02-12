package organization_user

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"

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
	data, err := copierx.Trans[types.OrganizationUserAssignRolesRequest, lender.OrganizationUserAssignRolesRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcReq := &lender.OrganizationUserAssignRolesRequest{
		Data: data,
	}

	_, err = l.svcCtx.OrganizationUserRpcClient.AssignRoles(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResponse{}, nil
}

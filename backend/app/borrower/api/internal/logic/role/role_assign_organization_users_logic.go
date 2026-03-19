package role

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleAssignOrganizationUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// assign organization users
func NewRoleAssignOrganizationUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleAssignOrganizationUsersLogic {
	return &RoleAssignOrganizationUsersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleAssignOrganizationUsersLogic) RoleAssignOrganizationUsers(req *types.RoleAssignOrganizationUserRequest) (resp *types.RoleEntity, err error) {
	rpcReq := &borrower.RoleAssignOrganizationUsersRequest{
		Data: &borrower.RoleAssignOrganizationUsersRequestData{},
	}

	rpcReq.Data, err = copierx.Trans[types.RoleAssignOrganizationUserRequest, borrower.RoleAssignOrganizationUsersRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.RoleRpcClient.AssignOrganizationUsers(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// trans response
	resp, err = copierx.Trans[borrower.RoleEntity, types.RoleEntity](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return resp, err
}

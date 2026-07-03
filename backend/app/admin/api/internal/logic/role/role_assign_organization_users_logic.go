package role

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"

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
	rpcReq := &admin.RoleAssignOrganizationUsersRequest{
		Data: &admin.RoleAssignOrganizationUsersRequestData{},
	}

	rpcReq.Data, err = copierx.Trans[types.RoleAssignOrganizationUserRequest, admin.RoleAssignOrganizationUsersRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.RoleRpcClient.AssignOrganizationUsers(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// trans response
	resp, err = copierx.Trans[proto.RoleEntity, types.RoleEntity](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return resp, err
}

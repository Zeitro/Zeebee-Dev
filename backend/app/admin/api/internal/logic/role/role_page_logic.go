package role

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
)

type RolePageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// page
func NewRolePageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RolePageLogic {
	return &RolePageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RolePageLogic) RolePage(req *types.RolePageRequest) (resp *types.RolePageResponse, err error) {
	resp = &types.RolePageResponse{}
	rpcReq := &admin.RolePageRequest{
		Data: &admin.RolePageRequestData{
			Filter: &admin.RolePageRequestFilter{},
		},
	}
	rpcReq.Data, err = copierx.Trans[types.RolePageRequest, admin.RolePageRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.RoleRpcClient.Page(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// trans response
	resp, err = copierx.Trans[admin.RolePageResponseData, types.RolePageResponse](rpcResp.Data)
	if err != nil {
		return nil, err
	}

	return resp, err
}

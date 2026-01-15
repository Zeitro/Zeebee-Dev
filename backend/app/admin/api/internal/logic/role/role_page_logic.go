package role

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"
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
	err = copier.CopyWithOption(rpcReq.Data, req, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}

	rpcResp, err := l.svcCtx.RoleRpcClient.Page(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// trans response
	err = copier.CopyWithOption(&resp, &rpcResp.Data, copier.Option{DeepCopy: true, CaseSensitive: true})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}
	return resp, err
}

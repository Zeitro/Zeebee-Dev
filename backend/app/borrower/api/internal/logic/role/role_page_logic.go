package role

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"
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
	rpcReq := &proto.CommonPageRequest{
		Data: &proto.CommonPageRequestData{
			Filter: &proto.CommonPageFilter{},
		},
	}
	rpcReq.Data, err = copierx.Trans[types.RolePageRequest, proto.CommonPageRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.RoleRpcClient.Page(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// trans response
	resp, err = copierx.Trans[borrower.RolePageResponseData, types.RolePageResponse](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return resp, err
}

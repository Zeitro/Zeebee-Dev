package role

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
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
	rpcReq := &proto.CommonPageRequest{
		Data: &proto.CommonPageRequestData{
			Filter: &proto.CommonPageFilter{},
		},
	}
	data, err := copierx.Trans[types.RolePageRequest, proto.CommonPageRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}
	rpcReq.Data = data

	rpcResp, err := l.svcCtx.RoleRpcClient.Page(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// trans response
	resp, err = copierx.Trans[lender.RolePageResponseData, types.RolePageResponse](rpcResp.Data, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}
	return resp, err
}

package role

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"
)

type RoleDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// detail
func NewRoleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleDetailLogic {
	return &RoleDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleDetailLogic) RoleDetail(req *types.PathIdRequest) (resp *types.RoleEntity, err error) {
	rpcReq := &proto.PathIdRequest{
		Data: &proto.PathIdRequestData{},
	}
	rpcReq.Data.Id = req.Id

	rpcResp, err := l.svcCtx.RoleRpcClient.Detail(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// trans response
	resp, err = copierx.Trans[proto.RoleEntity, types.RoleEntity](rpcResp.Data, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}
	return resp, err
}

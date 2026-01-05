package user

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// register
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.UserRegisterRequest) (resp *types.UserRegisterResponse, err error) {
	rpcReq := &lender.UserRegisterRequest{}
	rpcReq.Data, err = copierx.Trans[types.UserRegisterRequest, lender.UserRegisterRequestData](req)
	if err != nil {
		return nil, err
	}
	// call rpc server
	rpcResp, err := l.svcCtx.UserRpcClient.Register(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp, err = copierx.Trans[lender.UserRegisterResponseData, types.UserRegisterResponse](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return resp, err
}

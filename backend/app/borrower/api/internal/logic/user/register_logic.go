package user

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// register user
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.UserRegisterRequest) (resp *types.UserRegisterResponse, err error) {
	resp = &types.UserRegisterResponse{}
	rpcReq := &borrower.UserRegisterRequest{
		Data: &borrower.UserRegisterRequestData{},
	}

	rpcReq.Data, err = copierx.Trans[types.UserRegisterRequest, borrower.UserRegisterRequestData](req)
	if err != nil {
		return nil, err
	}

	// call rpc server
	rpcResp, err := l.svcCtx.UserRpcClient.Register(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// trans response data
	resp, err = copierx.Trans[borrower.UserRegisterResponseData, types.UserRegisterResponse](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return resp, err
}

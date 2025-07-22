package user

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
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
	rpcReq := &borrower.UserRegisterRequest{}
	resp = &types.UserRegisterResponse{}
	// trans request data
	err = copier.Copy(rpcReq, req)
	if err != nil {
		return nil, errors.Wrap(err, "trans request error")
	}
	// call rpc server
	rpcResp, err := l.svcCtx.UserRpcClient.Register(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// trans response data
	err = copier.Copy(resp, rpcResp)
	if err != nil {
		return nil, errors.Wrap(err, "trans response error")
	}

	return resp, err
}

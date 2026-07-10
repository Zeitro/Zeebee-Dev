package user

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"
)

type VerifyEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// verify email
func NewVerifyEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyEmailLogic {
	return &VerifyEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VerifyEmailLogic) VerifyEmail(req *types.UserVerifyEmailRequest) (resp *types.UserLoginResponse, err error) {
	resp = &types.UserLoginResponse{}

	rpcResp, err := l.svcCtx.UserRpcClient.VerifyEmail(l.ctx, &proto.VerifyEmailRequest{
		Data: &proto.VerifyEmailRequestData{
			Email: req.Email,
			Code:  req.Code,
		},
	})
	if err != nil {
		return nil, err
	}
	// trans response data
	resp, err = copierx.Trans[proto.UserLoginResponseData, types.UserLoginResponse](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

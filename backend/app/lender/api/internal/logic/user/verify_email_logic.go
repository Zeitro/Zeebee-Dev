package user

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"
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

func (l *VerifyEmailLogic) VerifyEmail(req *types.UserVerifyEmailRequest) (resp *types.UserJwtResponse, err error) {
	resp = &types.UserJwtResponse{}

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
	resp, err = copierx.Trans[proto.UserLoginResponseData, types.UserJwtResponse](rpcResp.Data)

	return resp, err
}

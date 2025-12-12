package user

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
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
	resp = &types.UserRegisterResponse{}
	// call rpc server
	rpcResp, err := l.svcCtx.UserRpcClient.Register(l.ctx, &lender.UserRegisterRequest{
		Data: &lender.UserRegisterRequestData{
			Email:          req.Email,
			Password:       req.Password,
			FirstName:      req.FirstName,
			LastName:       req.LastName,
			OrganizationId: req.OrganizationId,
		},
	})
	if err != nil {
		return nil, err
	}

	// trans response data
	err = copier.Copy(resp, &rpcResp.Data)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return resp, err
}

package user

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
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
	// call rpc server
	rpcResp, err := l.svcCtx.UserRpcClient.Register(l.ctx, &borrower.UserRegisterRequest{
		Data: &borrower.UserRegisterRequestData{
			Email:                 req.Email,
			Password:              req.Password,
			FirstName:             req.FirstName,
			LastName:              req.LastName,
			PreApprovalRequestId:  req.PreApprovalRequestId,
			TargetLenderCompanyId: req.TargetLenderCompanyId,

			// PhoneNumber: req.PhoneNumber,
			// Address:     req.Address,
			// Ssn: req.Ssn,
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

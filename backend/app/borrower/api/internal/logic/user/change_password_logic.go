package user

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// change password
func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangePasswordLogic) ChangePassword(req *types.UserChangePasswordRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &borrower.ChangePasswordRequest{
		Data: &borrower.ChangePasswordRequestData{},
	}
	err = copier.CopyWithOption(rpcReq.Data, req, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
	})
	if err != nil {
		err = errorx.GetFailedTransReqErr(err)
		return nil, err
	}

	_, err = l.svcCtx.UserRpcClient.ChangePassword(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	resp = &types.EmptyResponse{}
	return resp, nil
}

package user

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"
)

type EditProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// edit profile
func NewEditProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditProfileLogic {
	return &EditProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EditProfileLogic) EditProfile(req *types.UserEditProfileRequest) (resp *types.UserInfoResponse, err error) {
	resp = &types.UserInfoResponse{}
	rpcReq := &borrower.UserProfileEditRequest{
		Data: &borrower.UserProfileEditRequestData{},
	}

	rpcReq.Data, err = copierx.Trans[types.UserEditProfileRequest, borrower.UserProfileEditRequestData](req)

	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.UserRpcClient.EditProfile(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	resp, err = copierx.Trans[proto.UserinfoResponseData, types.UserInfoResponse](rpcResp.Data)
	if err != nil {
		return nil, err
	}

	return resp, err
}

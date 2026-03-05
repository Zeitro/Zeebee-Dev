package user

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
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

	// claim, err := utils.GetUserClaimFromCtx(l.ctx)
	// if err != nil {
	// 	return nil, err
	// }
	// userId := claim.UserId
	// if userId == 0 {
	// 	return nil, errorx.WithStack(errorx.New("User Id Not Found"))
	// }

	rpcReq := &admin.UserProfileEditRequest{
		Data: &admin.UserProfileEditRequestData{
			FirstName: req.FirstName,
			LastName:  req.LastName,
		},
	}

	// call rpc
	rpcResp, err := l.svcCtx.UserRpcClient.EditProfile(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// trans response
	err = copier.Copy(resp, &rpcResp.Data)
	resp, err = copierx.Trans[admin.UserinfoResponseData, types.UserInfoResponse](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return resp, err
}

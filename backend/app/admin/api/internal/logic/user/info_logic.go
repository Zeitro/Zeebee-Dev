package user

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
)

type InfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// info
func NewInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoLogic {
	return &InfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InfoLogic) Info(req *types.EmptyRequest) (resp *types.UserInfoResponse, err error) {
	resp = &types.UserInfoResponse{}

	// claim, err := utils.GetUserClaimFromCtx(l.ctx)
	// if err != nil {
	// 	return nil, err
	// }
	// userId := claim.UserId
	// if userId == 0 {
	// 	return nil, errorx.WithStack(errorx.New("User Id Not Found"))
	// }

	rpcResp, err := l.svcCtx.UserRpcClient.UserInfo(l.ctx, &admin.EmptyRequest{})
	if err != nil {
		return nil, err
	}

	// trans response
	resp, err = copierx.Trans[admin.UserinfoResponseData, types.UserInfoResponse](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return resp, err
}

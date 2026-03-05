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

type InfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// user info
func NewInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoLogic {
	return &InfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InfoLogic) Info(req *types.EmptyRequest) (resp *types.UserInfoResponse, err error) {
	resp = &types.UserInfoResponse{}
	userInfo, err := l.svcCtx.UserRpcClient.UserInfo(l.ctx, &borrower.EmptyRequest{})
	if err != nil {
		return nil, err
	}
	err = copier.Copy(resp, &userInfo.Data)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return resp, err
}

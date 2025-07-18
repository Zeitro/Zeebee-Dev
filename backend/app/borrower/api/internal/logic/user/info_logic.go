package user

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
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
	userId := utils.GetInt64FromCtx(l.ctx, constants.BorrowerUserIdKey)
	if userId == 0 {
		return nil, errors.New("user id not found")
	}
	userInfo, err := l.svcCtx.RpcClients.UserRpcClient.GetUserInfo(l.ctx, &borrower.UserInfoRequest{
		Id: userId,
	})
	if err != nil {
		return nil, err
	}
	err = copier.Copy(resp, userInfo)

	return resp, err
}

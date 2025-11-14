package user

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
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
	userId := utils.GetInt64FromCtx(l.ctx, constants.LenderUserIdKey)
	if userId == 0 {
		return nil, errors.WithStack(errors.New("User Id Not Found"))
	}
	userInfo, err := l.svcCtx.UserRpcClient.UserInfo(l.ctx, &lender.UserInfoRequest{
		Data: &lender.UserInfoRequestData{
			Id: userId,
		},
	})
	if err != nil {
		return nil, err
	}
	err = copier.Copy(resp, &userInfo.Data)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return resp, err
}

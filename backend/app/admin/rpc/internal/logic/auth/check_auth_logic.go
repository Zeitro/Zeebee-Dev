package auth

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/data/dto"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type CheckAuthLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewCheckAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckAuthLogic {
	return &CheckAuthLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *CheckAuthLogic) CheckAuth(in *admin.CheckAuthRequest) (resp *admin.CheckAuthResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	input, err := copierx.Trans[admin.CheckAuthRequestData, dto.AuthCheckRequestDto](in.Data)
	if err != nil {
		return nil, err
	}
	result, err := l.baseLogic.AuthLogic.AuthCheck(l.ctx, *input)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[dto.AuthCheckResponseDto, admin.CheckAuthResponseData](result)
	if err != nil {
		return nil, err
	}

	return &admin.CheckAuthResponse{
		Data: data,
	}, nil
}

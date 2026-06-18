package authlogic

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

type AuthCheckLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewAuthCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthCheckLogic {
	return &AuthCheckLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *AuthCheckLogic) AuthCheck(in *admin.AuthCheckRequest) (resp *admin.AuthCheckResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()

	input, err := copierx.Trans[admin.AuthCheckRequestData, dto.AuthCheckRequestDto](in.Data)
	if err != nil {
		return nil, err
	}
	result, err := l.baseLogic.AuthLogic.AuthCheck(l.ctx, *input)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[dto.AuthCheckResponseDto, admin.AuthCheckResponseData](result)
	if err != nil {
		return nil, err
	}

	return &admin.AuthCheckResponse{
		Data: data,
	}, nil
}

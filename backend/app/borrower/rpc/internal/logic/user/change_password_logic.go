package userlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/data/dto"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type ChangePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *ChangePasswordLogic) ChangePassword(in *borrower.ChangePasswordRequest) (resp *borrower.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	input := dto.UserChangePasswordRequestDto{
		OldPassword: in.Data.OldPassword,
		NewPassword: in.Data.NewPassword,
		UserId:      claim.UserId,
		UserType:    claim.UserType,
	}
	err = l.baseLogic.UserLogic.ChangePassword(l.ctx, input)
	if err != nil {
		return nil, err
	}

	return &borrower.EmptyResponse{}, nil
}

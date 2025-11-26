package userlogic

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/stringx"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type EditProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewEditProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditProfileLogic {
	return &EditProfileLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

func (l *EditProfileLogic) EditProfile(in *borrower.UserProfileEditRequest) (resp *borrower.UserinfoResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	user, err := l.svcCtx.UserRepo.GetUser(l.ctx, in.Data.UserId, constants.UserTypeBorrower)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	err = copier.CopyWithOption(user, in.Data, copier.Option{
		CaseSensitive: true,
	})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	user.FullName = stringx.FormatFullName(user.FirstName, user.MidName, user.LastName)
	user, err = l.svcCtx.UserRepo.UpdateOneById(l.ctx, user.Id, user)
	if err != nil {
		return nil, err
	}
	token, err := l.baseLogic.UserLogic.Info(l.ctx, user.Id, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}
	data := &borrower.UserinfoResponseData{}
	err = copier.CopyWithOption(data, token, copier.Option{
		CaseSensitive: true,
	})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}

	return &borrower.UserinfoResponse{
		Data: data,
	}, nil
}

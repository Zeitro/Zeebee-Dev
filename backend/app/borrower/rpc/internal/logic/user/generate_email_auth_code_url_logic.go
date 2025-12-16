package userlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type GenerateEmailAuthCodeUrlLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGenerateEmailAuthCodeUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateEmailAuthCodeUrlLogic {
	return &GenerateEmailAuthCodeUrlLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GenerateEmailAuthCodeUrlLogic) GenerateEmailAuthCodeUrl(in *borrower.UserEmailAuthCodeUrlRequest) (resp *borrower.UserEmailAuthCodeUrlResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	// todo: add your logic here and delete this line

	return &borrower.UserEmailAuthCodeUrlResponse{}, nil
}

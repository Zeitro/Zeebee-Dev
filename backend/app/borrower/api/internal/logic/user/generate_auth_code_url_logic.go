package user

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateAuthCodeUrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// generate auth code url
func NewGenerateAuthCodeUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateAuthCodeUrlLogic {
	return &GenerateAuthCodeUrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenerateAuthCodeUrlLogic) GenerateAuthCodeUrl(req *types.UserGenerateEmailAuthCodeUrlRequest) (resp *types.UserGenerateEmailAuthCodeUrlResponse, err error) {
	// todo: add your logic here and delete this line

	return
}

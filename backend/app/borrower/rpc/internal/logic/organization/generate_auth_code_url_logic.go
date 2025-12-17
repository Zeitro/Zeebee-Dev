package organizationlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type GenerateAuthCodeUrlLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGenerateAuthCodeUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateAuthCodeUrlLogic {
	return &GenerateAuthCodeUrlLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GenerateAuthCodeUrlLogic) GenerateAuthCodeUrl(in *borrower.OrganizationGenerateEmailCodeUrlRequest) (resp *borrower.OrganizationGenerateEmailCodeUrlResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	// todo: add your logic here and delete this line

	return &borrower.OrganizationGenerateEmailCodeUrlResponse{}, nil
}

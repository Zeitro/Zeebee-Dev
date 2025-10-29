package pdf

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyFormPdfLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// create apply form pdf
func NewApplyFormPdfLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyFormPdfLogic {
	return &ApplyFormPdfLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApplyFormPdfLogic) ApplyFormPdf(req *types.ApplyFormPdfRequest) (resp *types.EmptyResponse, err error) {
	// todo: add your logic here and delete this line

	return
}

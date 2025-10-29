package pdf

import (
	"bytes"
	"context"

	"zeebee-backend/app/borrower/api/internal/pdf/apply_form"
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

func (l *ApplyFormPdfLogic) ApplyFormPdf(req *types.ApplyFormPdfRequest) (contents []byte, err error) {

	var applyFormPdf *apply_form.ApplyForm
	var buf []byte
	ioWriter := bytes.NewBuffer(buf)
	applyFormPdf, err = apply_form.NewApplyForm(req)
	if err != nil {
		return
	}
	err = applyFormPdf.Generate(ioWriter)
	if err != nil {
		return
	}
	contents = ioWriter.Bytes()

	return
}

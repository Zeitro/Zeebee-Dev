package pdf_generator

import (
	"context"
	"strings"

	"zeebee-backend/app/pdf_generator/api/internal/global"
	"zeebee-backend/app/pdf_generator/api/internal/svc"
	"zeebee-backend/app/pdf_generator/api/internal/types"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type PdfCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewPdfCreateLogic creates a new PdfCreateLogic instance
func NewPdfCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PdfCreateLogic {
	return &PdfCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// PdfCreate generates a PDF based on the request type
func (l *PdfCreateLogic) PdfCreate(req *types.CreatePdfRequest) (resp *types.CreatePdfResponse, err error) {
	// Check if Chrome Pool is available
	if global.ChromePool == nil || !global.ChromePool.IsInitialized() {
		err = errorx.New("chrome pool is not available")
		return
	}

	pdfType := strings.ToLower(req.Type)

	// Generate PDF based on request type
	var pdfContent []byte
	ctx := context.WithoutCancel(l.ctx)
	pdfContent, err = global.ChromePool.GeneratePDF(ctx, pdfType, l.getContent(req))

	if err != nil {
		err = errorx.Wrap(err, "fail to create pdf")
		return
	}

	resp = &types.CreatePdfResponse{
		Data: pdfContent,
	}
	return
}

// getContent extracts the content based on request type
func (l *PdfCreateLogic) getContent(req *types.CreatePdfRequest) string {
	typ := strings.ToLower(req.Type)
	switch typ {
	case "url":
		return req.Url
	case "text":
		return req.Text
	default:
		return ""
	}
}

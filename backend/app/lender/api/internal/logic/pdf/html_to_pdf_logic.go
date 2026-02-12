package pdf

import (
	"context"
	"sync"

	"zeebee-backend/app/lender/api/internal/global"
	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/pdf_generator/client"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type HtmlToPdfLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// create apply form pdf
func NewHtmlToPdfLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HtmlToPdfLogic {
	return &HtmlToPdfLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HtmlToPdfLogic) HtmlToPdf(req *types.HtmlToPdfRequest) (result []byte, err error) {
	if req == nil {
		return nil, errorx.New("Request is empty")
	}

	if len(req.Content) == 0 {
		return nil, errorx.New("Content is empty")
	}

	cli := getPdfClient()
	result, err = cli.GeneratePdf(l.ctx, client.GeneratePdfRequest{
		Type: client.SourceTypeText,
		Text: req.Content,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

var pdfClient *client.Client
var pdfClientOnce sync.Once

func getPdfClient() *client.Client {
	pdfClientOnce.Do(func() {
		pdfClient = client.NewClient(global.Config.PdfGeneratorUrl)
	})
	return pdfClient
}

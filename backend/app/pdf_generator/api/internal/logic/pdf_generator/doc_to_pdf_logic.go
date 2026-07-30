package pdf_generator

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"zeebee-backend/app/pdf_generator/api/internal/libreoffice"
	"zeebee-backend/app/pdf_generator/api/internal/svc"
	"zeebee-backend/app/pdf_generator/api/internal/types"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type DocToPdfLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Convert doc/docx to Pdf
func NewDocToPdfLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DocToPdfLogic {
	return &DocToPdfLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// DocToPdf resolves the source bytes (Url download or inline Data), enforces the
// input cap, and converts via LibreOffice. Returned sentinels drive the HTTP
// status mapping in the handler.
func (l *DocToPdfLogic) DocToPdf(req *types.DocToPdfRequest) (resp *types.CreatePdfResponse, err error) {
	if l.svcCtx.Converter == nil {
		return nil, libreoffice.ErrSofficeNotFound
	}
	ext := strings.ToLower(req.Ext)

	src, label, err := l.sourceBytes(req, ext)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	pdf, err := l.svcCtx.Converter.Convert(l.ctx, src, ext)
	elapsed := time.Since(start).Milliseconds()
	l.Logger.WithContext(l.ctx).WithFields(
		logx.Field("biz", "doctopdf"),
		logx.Field("file", label),
		logx.Field("ext", ext),
		logx.Field("elapsed", elapsed),
		logx.Field("ok", err == nil),
	).Info("doc to pdf")
	if err != nil {
		return nil, err
	}
	return &types.CreatePdfResponse{Data: pdf}, nil
}

// sourceBytes fetches Url-mode bytes (capped) or returns Data-mode bytes inline.
func (l *DocToPdfLogic) sourceBytes(req *types.DocToPdfRequest, ext string) ([]byte, string, error) {
	switch ext {
	case "doc", "docx":
	default:
		return nil, "", errorx.Errorf("invalid file extension: %s", ext)
	}

	switch strings.ToLower(req.Type) {
	case "data":
		if len(req.Data) == 0 {
			return nil, "", libreoffice.ErrInvalidInput
		}
		if len(req.Data) > l.svcCtx.MaxInputBytes {
			return nil, "", libreoffice.ErrInvalidInput
		}
		return req.Data, "inline." + ext, nil
	default: // Url
		b, err := fetchURL(l.ctx, req.Url, l.svcCtx.MaxInputBytes)
		if err != nil {
			return nil, "", errorx.Wrap(err, "download source")
		}
		return b, redactURL(req.Url), nil
	}
}

// fetchURL GETs url with a byte cap matching MaxInputBytes.
func fetchURL(ctx context.Context, url string, maxBytes int) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("fail to download file %s with http status %d", redactURL(url), res.StatusCode)
	}
	b, err := io.ReadAll(io.LimitReader(res.Body, int64(maxBytes)+1))
	if err != nil {
		return nil, err
	}
	if len(b) > maxBytes {
		return nil, fmt.Errorf("source exceeds %d bytes", maxBytes)
	}
	return b, nil
}

// redactURL strips query and fragment (presigned S3 urls carry signing params
// that must not reach logs/errors). Keeps scheme+host+path for debug.
func redactURL(s string) string {
	u, err := url.Parse(s)
	if err != nil {
		return "<invalid url>"
	}
	u.RawQuery = ""
	u.Fragment = ""
	return u.String()
}

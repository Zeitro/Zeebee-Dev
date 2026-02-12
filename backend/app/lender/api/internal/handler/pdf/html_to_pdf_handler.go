package pdf

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/duke-git/lancet/v2/random"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"

	"zeebee-backend/app/lender/api/internal/logic/pdf"
	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/response"
	"zeebee-backend/common/validator"
)

// HtmlToPdfHandler convert html to pdf
func HtmlToPdfHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			if err != nil {
				logx.WithContext(r.Context()).
					WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).
					Error(err)
			}
		}()
		var req types.HtmlToPdfRequest
		if err = httpx.Parse(r, &req); err != nil {
			response.BadRequest(w, errorx.GetErrorMsg(err))
			return
		}

		v, err := validator.GetEnValidator()
		if err != nil {
			response.InternalServerError(w, errorx.GetErrorMsg(err))
			return
		}
		if err := v.Validate(&req); err != nil {
			response.BadRequest(w, errorx.GetErrorMsg(err))
			return
		}

		l := pdf.NewHtmlToPdfLogic(r.Context(), svcCtx)
		resp, err := l.HtmlToPdf(&req)
		if err != nil {
			response.BadRequest(w, errorx.GetErrorMsg(err))
			return
		}
		fileName := req.Name
		if len(fileName) == 0 {
			fileName = random.RandString(10)
		}
		fileName = strings.ToLower(fileName)
		if !strings.HasSuffix(fileName, ".pdf") {
			fileName = fileName + ".pdf"
		}

		// output pdf
		w.Header().Set("Content-Type", "application/pdf")
		// w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"; filename*=utf-8''%s", fileName, url.PathEscape(fileName)))
		w.Header().Set("Content-Length", strconv.FormatInt(int64(len(resp)), 10))
		_, _ = w.Write(resp)
	}
}

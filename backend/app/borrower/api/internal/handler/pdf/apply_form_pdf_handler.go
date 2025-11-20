package pdf

import (
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"

	"zeebee-backend/app/borrower/api/internal/logic/pdf"
	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/response"
	"zeebee-backend/common/validator"
)

// create apply form pdf
func ApplyFormPdfHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			if err != nil {
				logx.WithContext(r.Context()).
					WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).
					Error(err)
			}
		}()
		var req types.ApplyFormPdfRequest
		var v *validator.Validator
		// var applyFormPdf *apply_form.ApplyForm

		if err = httpx.Parse(r, &req); err != nil {
			// response.BadRequestWithCode(w, http.StatusBadRequest, errorx.GetErrorMsg(err))
			response.BadRequest(w, errorx.GetErrorMsg(err))
			return
		}

		v, err = validator.GetEnValidator()
		if err != nil {
			// response.BadRequestWithCode(w, http.StatusBadRequest, errorx.GetErrorMsg(err))
			response.BadRequest(w, errorx.GetErrorMsg(err))
			return
		}
		if err = v.Validate(&req); err != nil {
			// response.BadRequestWithCode(w, http.StatusBadRequest, errorx.GetErrorMsg(err))
			response.BadRequest(w, errorx.GetErrorMsg(err))
			return
		}

		l := pdf.NewApplyFormPdfLogic(r.Context(), svcCtx)
		resp, err := l.ApplyFormPdf(&req)

		// applyFormPdf, err = apply_form.NewApplyForm(&req)
		// if err != nil {
		// 	// response.BadRequestWithCode(w, http.StatusBadRequest, errorx.GetErrorMsg(err))
		// 	response.BadRequest(w, errorx.GetErrorMsg(err))
		// 	return
		// }
		//
		// var buf []byte
		// ioWriter := bytes.NewBuffer(buf)
		//
		// err = applyFormPdf.Generate(ioWriter)
		if err != nil {
			// response.BadRequestWithCode(w, http.StatusBadRequest, errorx.GetErrorMsg(err))
			response.BadRequest(w, errorx.GetErrorMsg(err))
			return
		}

		// output pdf
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "attachment; filename=apply_form.pdf")
		w.Header().Set("Content-Length", strconv.FormatInt(int64(len(resp)), 10))
		_, _ = w.Write(resp)
	}
}

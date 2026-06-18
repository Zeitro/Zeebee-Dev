package loan

import (
	"io"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"

	"zeebee-backend/app/borrower/api/internal/logic/loan"
	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/response"
)

// ImportHandler import
func ImportHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			if err != nil {
				logx.WithContext(r.Context()).
					WithFields(
						logx.Field("stacks", errorx.GetErrorStack(err)),
						logx.Field("service_name", svcCtx.Config.Name),
					).
					Error(err)
			}
		}()

		var req types.ImportLoanRequest
		if err = httpx.Parse(r, &req); err != nil {
			response.BadRequest(w, errorx.GetErrorMsg(err))
			return
		}

		_ = r.ParseMultipartForm(32 << 20)
		f, _, err := r.FormFile("file")
		if err != nil {
			response.BadRequest(w, errorx.GetErrorMsg(err))
			return
		}
		defer f.Close()
		var contents []byte
		contents, err = io.ReadAll(f)
		if err != nil {
			response.BadRequest(w, errorx.GetErrorMsg(err))
			return
		}

		l := loan.NewImportLogic(r.Context(), svcCtx)
		resp, err := l.Import(contents, req.Stage)
		if err != nil {
			response.BadRequest(w, errorx.GetErrorMsg(err))
		} else {
			response.Success(w, resp)
		}
	}
}

package loan

import (
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/logic/loan"
	"zeebee-backend/app/borrower/api/internal/svc"
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
					WithFields(logx.Field("stacks", errorx.GetErrorStack(errors.WithStack(err)))).
					Error(err)
			}
		}()
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
		resp, err := l.Import(contents)
		if err != nil {
			response.BadRequest(w, errorx.GetErrorMsg(err))
		} else {
			response.Success(w, resp)
		}
	}
}

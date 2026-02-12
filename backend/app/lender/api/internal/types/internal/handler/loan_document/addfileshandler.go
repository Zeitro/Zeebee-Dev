// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package loan_document

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"zeebee-backend/app/lender/api/internal/types/internal/logic/loan_document"
	"zeebee-backend/app/lender/api/internal/types/internal/svc"
	"zeebee-backend/app/lender/api/internal/types/internal/types"
)

// add files
func AddFilesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoanDocumentAddFilesRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := loan_document.NewAddFilesLogic(r.Context(), svcCtx)
		resp, err := l.AddFiles(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package loan

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"zeebee-backend/app/lender/api/internal/types/internal/logic/loan"
	"zeebee-backend/app/lender/api/internal/types/internal/svc"
	"zeebee-backend/app/lender/api/internal/types/internal/types"
)

// page
func PageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoanPageRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := loan.NewPageLogic(r.Context(), svcCtx)
		resp, err := l.Page(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

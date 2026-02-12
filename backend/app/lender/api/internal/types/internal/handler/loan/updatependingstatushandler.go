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

// update loan pending status
func UpdatePendingStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateLoanPendingStatusRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := loan.NewUpdatePendingStatusLogic(r.Context(), svcCtx)
		resp, err := l.UpdatePendingStatus(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

package loan

import (
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"

	"zeebee-backend/app/borrower/api/internal/logic/loan"
	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/response"
)

// create draft
func CreateDraftHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
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
		var req types.LoanApplyRequest
		if err = httpx.Parse(r, &req); err != nil {
			response.BadRequest(w, errorx.GetErrorMsg(err))
			return
		}
		l := loan.NewDraftLogic(r.Context(), svcCtx)
		resp, err := l.CreateDraft(&req)
		if err != nil {
			response.BadRequest(w, errorx.GetErrorMsg(err))
		} else {
			response.Success(w, resp)
		}
	}
}

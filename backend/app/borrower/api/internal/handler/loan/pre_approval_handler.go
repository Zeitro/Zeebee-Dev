package loan

import (
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"

	"zeebee-backend/app/borrower/api/internal/logic/loan"
	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/api/internal/validate_rules"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/response"
	"zeebee-backend/common/validator"
)

// pre_approve
func PreApprovalHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
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
		var req types.LoanPreApprovalRequest
		if err = httpx.Parse(r, &req); err != nil {
			response.BadRequest(w, errorx.GetErrorMsg(err))
			return
		}

		v, err := validator.GetEnValidator()
		if err != nil {
			response.InternalServerError(w, errorx.GetErrorMsg(err))
			return
		}
		v.RegisterStructRule(validate_rules.LoanPreApprovalRequestProperty, types.LoanPreApprovalRequestProperty{})
		_ = v.RegisterTansMsg(validate_rules.PropertyUnitTag, validate_rules.PropertyUnitMsgTpl)
		if err := v.Validate(&req); err != nil {
			response.BadRequest(w, errorx.GetErrorMsg(err))
			return
		}

		l := loan.NewPreApprovalLogic(r.Context(), svcCtx)
		resp, err := l.PreApproval(&req)
		if err != nil {
			response.BadRequest(w, errorx.GetErrorMsg(err))
		} else {
			response.Success(w, resp)
		}
	}
}

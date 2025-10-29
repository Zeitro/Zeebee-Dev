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
	"zeebee-backend/common/validator"
)

// update apply
func UpdateApplyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			if err != nil {
				logx.WithContext(r.Context()).
					WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).
					Error(err)
			}
		}()
		var req types.LoanApplyRequest
		if err = httpx.Parse(r, &req); err != nil {
			response.BadRequest(w, errorx.GetErrorMsg(err))
			return
		}
		// if req.LoanProperty != nil && req.LoanProperty.PropertyBase != nil {
		// 	req.LoanProperty.Unit = req.LoanProperty.PropertyBase.Unit
		// }

		v, err := validator.GetEnValidator()
		if err != nil {
			response.InternalServerError(w, errorx.GetErrorMsg(err))
			return
		}
		// v.RegisterStructRule(validate_rules.LoanApplyPropertyEntry, types.LoanApplyPropertyEntry{})
		// _ = v.RegisterTansMsg(validate_rules.PropertyUnitTag, validate_rules.PropertyUnitMsgTpl)
		if err := v.Validate(&req); err != nil {
			response.BadRequest(w, errorx.GetErrorMsg(err))
			return
		}

		l := loan.NewUpdateApplyLogic(r.Context(), svcCtx)
		resp, err := l.UpdateApply(&req)
		if err != nil {
			response.BadRequest(w, errorx.GetErrorMsg(err))
		} else {
			response.Success(w, resp)
		}
	}
}

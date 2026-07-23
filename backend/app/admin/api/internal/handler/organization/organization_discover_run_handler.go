package organization

import (
	"net/http"

	"zeebee-backend/app/admin/api/internal/logic/organization"
	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/response"
	"zeebee-backend/common/validator"
)

// lender discover run (scheduled trigger only; manual flow uses preview + commit)
func OrganizationDiscoverRunHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			if err != nil {
				svcCtx.Logger.WithContext(r.Context()).
					WithFields(
						logx.Field("stacks", errorx.GetErrorStack(err)),
					).
					Error(err)
			}
		}()
		var req types.LenderDiscoverRunRequest
		if err = httpx.Parse(r, &req); err != nil {
			response.BadRequest(w, errorx.GetErrorMsg(err))
			return
		}
		var v *validator.Validator
		v, err = validator.GetEnValidator()
		if err != nil {
			response.InternalServerError(w, errorx.GetErrorMsg(err))
			return
		}
		err = v.Validate(&req)
		if err != nil {
			response.BadRequest(w, errorx.GetErrorMsg(err))
			return
		}

		l := organization.NewOrganizationDiscoverRunLogic(r.Context(), svcCtx)
		resp, err := l.LenderDiscoverRun(&req)
		if err != nil {
			response.BadRequest(w, errorx.GetErrorMsg(err))
		} else {
			response.Success(w, resp)
		}
	}
}

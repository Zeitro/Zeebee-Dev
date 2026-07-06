package organization

import (
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"

	"zeebee-backend/app/admin/api/internal/logic/organization"
	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/response"
	"zeebee-backend/common/validator"
)

// switch status
func OrganizationSwitchStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
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
		var req types.OrganizationSwitchStatusRequest
		if err = httpx.Parse(r, &req); err != nil {
			response.BadRequest(w, errorx.GetErrorMsg(err))
			return
		}

		v, err := validator.GetEnValidator()
		if err != nil {
			response.InternalServerError(w, errorx.GetErrorMsg(err))
			return
		}
		if err := v.Validate(&req); err != nil {
			response.BadRequest(w, errorx.GetErrorMsg(err))
			return
		}

		l := organization.NewOrganizationSwitchStatusLogic(r.Context(), svcCtx)
		resp, err := l.OrganizationSwitchStatus(&req)
		if err != nil {
			response.BadRequest(w, errorx.GetErrorMsg(err))
		} else {
			response.Success(w, resp)
		}
	}
}

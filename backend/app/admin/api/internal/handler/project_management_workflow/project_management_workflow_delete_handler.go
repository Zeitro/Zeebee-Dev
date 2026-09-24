package project_management_workflow

import (
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"zeebee-backend/app/admin/api/internal/logic/project_management_workflow"
	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/response"
	"zeebee-backend/common/validator"
)

// delete
func ProjectManagementWorkflowDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
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
		var req types.ProjectManagementWorkflowDetailRequest
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

		l := project_management_workflow.NewProjectManagementWorkflowDeleteLogic(r.Context(), svcCtx)
		resp, err := l.ProjectManagementWorkflowDelete(&req)
		if err != nil {
			response.BadRequest(w, errorx.GetErrorMsg(err))
		} else {
			response.Success(w, resp)
		}
	}
}

// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package organization

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"zeebee-backend/app/lender/api/internal/types/internal/logic/organization"
	"zeebee-backend/app/lender/api/internal/types/internal/svc"
	"zeebee-backend/app/lender/api/internal/types/internal/types"
)

// update profile
func OrganizationUpdateProfileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OrganizationUpdateProfileRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := organization.NewOrganizationUpdateProfileLogic(r.Context(), svcCtx)
		resp, err := l.OrganizationUpdateProfile(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

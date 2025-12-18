// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package organization

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"zeebee-backend/app/borrower/api/internal/logic/organization"
	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
)

// generate auth code url
func OrganizationGenerateAuthCodeUrlHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OrganizationGenerateAuthCodeUrlRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := organization.NewOrganizationGenerateAuthCodeUrlLogic(r.Context(), svcCtx)
		resp, err := l.OrganizationGenerateAuthCodeUrl(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package organization_email

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"zeebee-backend/app/borrower/api/internal/logic/organization_email"
	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
)

// generate auth code url
func OrganizationEmailGenerateAuthCodeUrlHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OrganizationEmailGenerateAuthCodeUrlRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := organization_email.NewOrganizationEmailGenerateAuthCodeUrlLogic(r.Context(), svcCtx)
		resp, err := l.OrganizationEmailGenerateAuthCodeUrl(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"zeebee-backend/app/lender/api/internal/types/internal/logic/user"
	"zeebee-backend/app/lender/api/internal/types/internal/svc"
	"zeebee-backend/app/lender/api/internal/types/internal/types"
)

// send verify email
func SendVerifyEmailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserSendVerifyEmailRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewSendVerifyEmailLogic(r.Context(), svcCtx)
		resp, err := l.SendVerifyEmail(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

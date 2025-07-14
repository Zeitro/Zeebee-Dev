package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"zeebee-backend/app/borrower/api/internal/logic"
	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"

	"zeebee-backend/common/response"
)

func GetUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			response.BadRequest(w, err.Error())
			return
		}

		l := logic.NewGetUserLogic(r.Context(), svcCtx)
		resp, err := l.GetUser(&req)
		if err != nil {
			response.Success(w, resp)
		} else {
			response.BadRequest(w, err.Error())
		}
	}
}

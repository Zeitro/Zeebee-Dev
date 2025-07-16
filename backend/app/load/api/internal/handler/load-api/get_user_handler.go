package load-api

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"zeebee-backend/app/load/api/internal/logic/load-api"
	"zeebee-backend/app/load/api/internal/svc"
	"zeebee-backend/app/load/api/internal/types"


	xhttp "github.com/zeromicro/x/http"
)


func GetUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
			return
		}

		l := load-api.NewGetUserLogic(r.Context(), svcCtx)
		resp, err := l.GetUser(&req)
		if err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
		} else {
			xhttp.JsonBaseResponseCtx(r.Context(), w, resp)
		}
	}
}

// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package storage

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"zeebee-backend/app/admin/api/internal/types/internal/logic/storage"
	"zeebee-backend/app/admin/api/internal/types/internal/svc"
	"zeebee-backend/app/admin/api/internal/types/internal/types"
)

// create aws s3 presign url
func CreateAwsS3PresignUrlHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.S3PresignUrlCreateRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := storage.NewCreateAwsS3PresignUrlLogic(r.Context(), svcCtx)
		resp, err := l.CreateAwsS3PresignUrl(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

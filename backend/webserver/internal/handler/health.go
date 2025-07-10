package handler

import (
	"net/http"

	"zeebee-backend/common/response"
	"zeebee-backend/webserver/internal/svc"
)

type HealthHandler struct {
	*BaseHandler
}

func NewHealthHandler(svcCtx *svc.ServiceContext) *HealthHandler {
	return &HealthHandler{
		BaseHandler: NewBaseHandler(svcCtx),
	}
}

// Health 健康检查
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]interface{}{
		"status":  "ok",
		"service": "webserver",
		"version": "1.0.0",
	})
}

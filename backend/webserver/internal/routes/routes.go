package routes

import (
	"net/http"

	"zeebee-backend/webserver/internal/handler"
	"zeebee-backend/webserver/internal/svc"
)

func RegisterRoutes(svcCtx *svc.ServiceContext) {
	// 创建处理器
	healthHandler := handler.NewHealthHandler(svcCtx)

	http.HandleFunc("/health", healthHandler.Health)

	// TODO: 添加更多路由

}

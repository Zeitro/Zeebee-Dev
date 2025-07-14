package routes

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest"

	"zeebee-backend/webserver/internal/handler"
	"zeebee-backend/webserver/internal/svc"
)

func RegisterRoutes(server *rest.Server, svcCtx *svc.ServiceContext) {
	// 创建处理器
	healthHandler := handler.NewHealthHandler(svcCtx)

	// http.HandleFunc("/health", healthHandler.Health)

	server.AddRoutes([]rest.Route{
		{
			Method:  http.MethodGet,
			Path:    "/health",
			Handler: healthHandler.Health,
		},
		// TODO: 添加更多路由
	})

}

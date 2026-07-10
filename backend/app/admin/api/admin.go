package api

import (
	_ "embed"

	"github.com/zeromicro/go-zero/core/service"

	"zeebee-backend/app/admin/api/internal/handler"
	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/config"
	"zeebee-backend/common/response"

	"github.com/zeromicro/go-zero/rest"
)

//go:embed admin.api
var adminApiContent string

type ApiServer struct {
	s *rest.Server
}

func NewApiServer(c config.Config) (s service.Service) {
	apiConfig := c.ApiConf
	runOpts := []rest.RunOption{
		rest.WithNotFoundHandler(response.NotFoundHandler{}),
		rest.WithNotAllowedHandler(response.NotFoundHandler{}),
		rest.WithCors("*"),
	}
	restServer := rest.MustNewServer(apiConfig.RestConf, runOpts...)
	// restServer.Use(middleware.CorsMiddleware)

	ctx := svc.NewServiceContext(apiConfig, adminApiContent)

	// register handlers
	handler.RegisterHandlers(restServer, ctx)

	// extra route
	handler.RegisterExtraRoutes(restServer)
	handler.RegisterStripeWebhookRoute(restServer, ctx)

	return &ApiServer{
		s: restServer,
	}
}

func (s *ApiServer) Start() {
	go func() {
		s.s.Start()
	}()
}

func (s *ApiServer) Stop() {
	s.s.Stop()
}

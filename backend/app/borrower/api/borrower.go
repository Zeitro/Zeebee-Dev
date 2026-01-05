package api

import (
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest"

	"zeebee-backend/app/borrower/api/internal/handler"
	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/config"
	"zeebee-backend/common/response"
)

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

	ctx := svc.NewServiceContext(apiConfig)

	// register handlers
	handler.RegisterHandlers(restServer, ctx)

	// extra route
	handler.RegisterExtraRoutes(restServer)

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

package api

import (
	_ "embed"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
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
	s        *rest.Server
	startOpt rest.StartOption
}

func NewApiServer(c config.Config) (s service.Service) {
	apiConfig := c.ApiConf
	runOpts := []rest.RunOption{
		rest.WithNotFoundHandler(response.NotFoundHandler{}),
		rest.WithNotAllowedHandler(response.NotAllowedHandler{}),
		rest.WithCors("*"),
	}
	restServer := rest.MustNewServer(apiConfig.RestConf, runOpts...)

	// Wrap the http.Server handler ABOVE go-zero's TraceHandler/LogHandler so the
	// built-in access log line carries service_name. Inner Use() middlewares run
	// below LogHandler and can't reach it.
	startOpt := func(svr *http.Server) {
		base := svr.Handler
		svr.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := logx.ContextWithFields(r.Context(),
				logx.Field("service_name", apiConfig.Name))
			base.ServeHTTP(w, r.WithContext(ctx))
		})
	}
	logger := logx.WithCallerSkip(0).
		WithFields(logx.Field("service_name", apiConfig.Name))

	ctx := svc.NewServiceContext(apiConfig, logger, adminApiContent)

	// register handlers
	handler.RegisterHandlers(restServer, ctx)

	// extra route
	handler.RegisterExtraRoutes(restServer)

	// public Stripe webhook endpoint (hand-written: needs raw body for
	// signature verification and bypasses AuthMiddleware)
	handler.RegisterStripeWebhookRoute(restServer, ctx)

	return &ApiServer{
		s:        restServer,
		startOpt: startOpt,
	}
}

func (s *ApiServer) Start() {
	go func() {
		s.s.StartWithOpts(s.startOpt)
	}()
}

func (s *ApiServer) Stop() {
	s.s.Stop()
}

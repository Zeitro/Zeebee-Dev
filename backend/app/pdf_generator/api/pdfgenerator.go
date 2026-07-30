package api

import (
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"

	"zeebee-backend/app/pdf_generator/api/internal/global"
	"zeebee-backend/app/pdf_generator/api/internal/handler"
	"zeebee-backend/app/pdf_generator/api/internal/pdf_generator"
	"zeebee-backend/app/pdf_generator/api/internal/svc"
	"zeebee-backend/app/pdf_generator/config"
	"zeebee-backend/common/response"

	"github.com/zeromicro/go-zero/rest"
)

type ApiServer struct {
	s        *rest.Server
	svcCtx   *svc.ServiceContext
	startOpt rest.StartOption
}

func NewApiServer(c config.Config) (s service.Service, err error) {
	// Initialize Chrome Pool during service startup
	// If initialization fails, the service cannot start
	maxConcurrent := pdf_generator.DefaultMaxConcurrent
	acquireTimeout := pdf_generator.DefaultAcquireTimeout
	if err = global.InitChromePool(maxConcurrent, acquireTimeout); err != nil {
		return nil, err
	}

	apiConfig := c.ApiConf
	runOpts := []rest.RunOption{
		rest.WithNotFoundHandler(response.NotFoundHandler{}),
		rest.WithNotAllowedHandler(response.NotAllowedHandler{}),
		rest.WithCors("*"),
	}
	restServer := rest.MustNewServer(apiConfig.RestConf, runOpts...)
	// restServer.Use(middleware.CorsMiddleware)

	// Wrap the http.Server handler ABOVE go-zero's TraceHandler/LogHandler so the
	// built-in access log line carries service_name.
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

	ctx := svc.NewServiceContext(apiConfig, logger)

	// register handlers
	handler.RegisterHandlers(restServer, ctx)

	// extra route
	handler.RegisterExtraRoutes(restServer)

	return &ApiServer{
		s:        restServer,
		svcCtx:   ctx,
		startOpt: startOpt,
	}, nil
}

func (s *ApiServer) Start() {
	go func() {
		s.s.StartWithOpts(s.startOpt)
	}()
}

func (s *ApiServer) Stop() {
	// Cleanup Chrome Pool resources
	if s.svcCtx != nil {
		s.svcCtx.Close()
	}

	s.s.Stop()
}
